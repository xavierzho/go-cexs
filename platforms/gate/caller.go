package gate

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (c *Connector) Sign(params []byte) string {
	mac := hmac.New(sha512.New, []byte(c.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

// Call Reference https://www.gate.io/docs/developers/apiv4/#apiv4-signed-request-requirements
func (c *Connector) Call(method string, route string, params platforms.Serializer, authType constants.AuthType, returnType interface{}) error {
	// Add necessary parameters
	symbol, ok := params.Exists(SymbolFiled)
	if ok {
		params.Set(SymbolFiled, c.SymbolPattern(symbol.(string)))
	}
	bytesBody, err := utils.Json.Marshal(params)
	if err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sha := sha512.New()
	sha.Write(bytesBody)
	reqBody := bytes.NewReader(bytesBody)
	queryString := ""
	switch method {
	case http.MethodDelete, http.MethodGet:
		query, err := params.EncodeQuery()
		if err != nil {
			return err
		}
		queryString = query
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s?%s", RestAPI, route, queryString), reqBody)
	if err != nil {
		return err
	}
	//fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, route, "", hex.EncodeToString(sha.Sum(nil)), timestamp)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add(TimestampHeader, timestamp)

	switch authType {
	case constants.Keyed:
		// key only request
		req.Header.Add(KeyHeader, c.APIKey)
	case constants.Signed:
		// must sign
		// prepare sign
		var signBytes = new(bytes.Buffer)
		signBytes.WriteString(method)
		signBytes.WriteByte('\n')
		signBytes.WriteString(route)
		signBytes.WriteByte('\n')
		signBytes.WriteString(queryString)
		signBytes.WriteByte('\n')
		signBytes.WriteString(hex.EncodeToString(sha.Sum(nil)))
		signBytes.WriteByte('\n')
		signBytes.WriteString(timestamp)
		signature := c.Sign(signBytes.Bytes())
		req.Header.Add(SignHeader, signature)
		req.Header.Add(KeyHeader, c.APIKey)
	default:
		// default None
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("")
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return utils.Json.Unmarshal(respBody, returnType)
}
