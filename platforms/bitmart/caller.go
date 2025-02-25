package bitmart

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/xavierzho/go-cexs/constants"
)

func (c *Connector) Sign(message []byte) string {
	var key = bytes.NewBufferString(c.APISecret)

	mac := hmac.New(sha256.New, key.Bytes())
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
	//req.Header.Add("X-BM-SIGN", hex.EncodeToString(mac.Sum(message)))
}
func preSign(timestamp time.Time, memo string, body []byte) []byte {
	var buf = bytes.NewBufferString(fmt.Sprintf("%d#%s#", timestamp.UnixNano(), memo))
	buf.Write(body)
	return buf.Bytes()
}

type Response struct {
	Code    int             `json:"code"`
	Trace   string          `json:"trace"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (c *Connector) Call(method string, route string, params platforms.Serializer, authType constants.AuthType,
	returnType any) error {
	var err error
	timestamp := time.Now()
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("User-Agent", "bitmart-python-sdk-api/")
	symbol, ok := params.Exists(SymbolFiled)
	if ok {
		params.Set(SymbolFiled, c.SymbolPattern(symbol.(string)))
	}
	body, err := params.Serialize()
	if err != nil {
		return err
	}
	var bodyData = new(bytes.Buffer)
	_, err = io.Copy(bodyData, body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%s", RestAPI, route)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	switch authType {
	case constants.Keyed:
		header.Add("X-BM-KEY", c.APIKey)
	case constants.Signed:
		header.Add("X-BM-KEY", c.APIKey)
		header.Add("X-BM-TIMESTAMP", strconv.FormatInt(timestamp.UnixNano(), 10))
		signature := c.Sign(preSign(timestamp, *c.Option, bodyData.Bytes()))
		req.Header.Add("X-BM-SIGN", signature)
	default:
		header.Add("X-BM-TIMESTAMP", strconv.FormatInt(timestamp.UnixNano(), 10))

	}

	if method == http.MethodGet {
		query, err := params.EncodeQuery()
		if err != nil {
			return err
		}
		url = fmt.Sprintf("%s?%s", url, query)
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else if method == http.MethodPost {
		req.Body = io.NopCloser(body)
	}
	var response Response
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Code != 1000 || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Bitmart](error code=%d) %s", response.Code, response.Message)
	}
	return json.Unmarshal(response.Data, returnType)
}
