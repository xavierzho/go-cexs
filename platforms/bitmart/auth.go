package bitmart

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/utils"
)

func (c *Connector) Sign(req *http.Request, message []byte) {
	var key = bytes.NewBufferString(c.APISecret)

	mac := hmac.New(sha256.New, key.Bytes())
	req.Header.Add("X-BM-SIGN", hex.EncodeToString(mac.Sum(message)))
}
func preSign(timestamp time.Time, memo string, body []byte) []byte {
	var buf = bytes.NewBufferString(fmt.Sprintf("%d#%s#", timestamp.UnixNano(), memo))
	buf.Write(body)
	return buf.Bytes()
}

type Response struct {
	Code    int         `json:"code"`
	Trace   string      `json:"trace"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *Connector) Call(method string, route string, params map[string]interface{}, authType constants.AuthType,
	returnType interface{}) error {
	timestamp := time.Now()
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("User-Agent", "bitmart-python-sdk-api/")
	bodyData, _ := json.Marshal(params)
	var body = bytes.NewBuffer(bodyData)
	url := fmt.Sprintf("%s/%s", RestAPI, route)

	var req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	if authType == constants.Keyed {
		header.Add("X-BM-KEY", c.APIKey)
	} else {
		header.Add("X-BM-KEY", c.APIKey)
		header.Add("X-BM-TIMESTAMP", strconv.FormatInt(timestamp.UnixNano(), 10))
		c.Sign(req, preSign(timestamp, *c.Option, bodyData))
	}

	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else if method == http.MethodPost {
		req.Body = io.NopCloser(body)
	}
	var response Response
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	err = utils.Json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Code != 1000 || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Bitmart](error code=%d) %s", response.Code, response.Message)
	}
	dataBytes, err := utils.Json.Marshal(response.Data)
	if err != nil {
		return err
	}
	err = utils.Json.Unmarshal(dataBytes, returnType)
	if err != nil {
		return err
	}
	return nil
}
