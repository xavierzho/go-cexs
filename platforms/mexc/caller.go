package mexc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"io"
	"net/http"
	"time"
)

func (c *Connector) Sign(params []byte) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Connector) Call(method string, route string, params platforms.Serializer, authType constants.AuthType, returnType interface{}) error {
	// Add necessary parameters
	var timestamp = time.Now()
	var url = RestAPI + route
	params.Set("timestamp", timestamp.UnixMilli())
	params.Set("recvWindow", 6000)
	//params["timestamp"] = timestamp.UnixMilli()
	//params["recvWindow"] = 6000
	symbol, ok := params.Exists(SymbolFiled)
	if ok {
		//params[SymbolFiled] = c.SymbolPattern(symbol.(string))
		params.Set(SymbolFiled, c.SymbolPattern(symbol.(string)))
	}
	bytesBody, err := json.Marshal(params)
	if err != nil {
		return err
	}
	header := new(http.Header)
	queryString, err := params.EncodeQuery()
	if err != nil {
		return err
	}
	switch method {
	case http.MethodGet:
		url += "?" + queryString
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		header.Add("Content-Type", "application/json")
	}
	switch authType {
	case constants.Keyed:
		// key only request
		header.Add(KeyHeader, c.APIKey)
	case constants.Signed:
		// must sign
		header.Add(KeyHeader, c.APIKey)
		signature := c.Sign([]byte(queryString))
		url = fmt.Sprintf("%s&signature=%s", url, signature)
	default:
		// default None
	}

	req, err := http.NewRequest(method, RestAPI+route+"?"+queryString, bytes.NewReader(bytesBody))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, returnType)
}
