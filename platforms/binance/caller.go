package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/utils"
)

func (c *Connector) Sign(params []byte) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Connector) Call(method string, route string, params platforms.Serializer, authType constants.AuthType, returnType interface{}) error {
	headers := http.Header{}
	var reqBody io.Reader = nil
	params.Set(TimeFiled, time.Now().UnixMilli())
	encoded, err := params.EncodeQuery()
	if err != nil {
		return err
	}

	var fullUrl = fmt.Sprintf("%s%s?%s", RestAPI, route, encoded)
	symbol, ok := params.Exists(SymbolFiled)
	if ok {
		params.Set(SymbolFiled, c.SymbolPattern(symbol.(string)))
	}
	switch authType {
	case constants.Keyed:
		headers.Add(HeaderAPIKEY, c.APIKey)
	case constants.Signed:
		headers.Add(HeaderAPIKEY, c.APIKey)
		signature := c.Sign([]byte(encoded))

		fullUrl += fmt.Sprintf("&%s=%s", SignatureFiled, url.QueryEscape(signature))
	default:
		// default none
	}

	fmt.Println("full url", fullUrl)
	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return err
	}

	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	decoder := utils.Json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err := decoder.Decode(&errResp)
		if err != nil {
			return err
		}
		return fmt.Errorf("binance request error(%d)[%s]", errResp.Code, errResp.Msg)
	}

	return decoder.Decode(&returnType)
}
