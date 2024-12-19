package binance

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/utils"
)

func EncodeParams(m map[string]any) string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result []string
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s=%v", key, m[key]))
	}
	return strings.Join(result, "&")
}

func (c *Connector) Sign(params []byte) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Connector) Call(method string, route string, params map[string]any, authType constants.AuthType, returnType interface{}) error {
	headers := http.Header{}
	var reqBody io.Reader = nil
	params[TimeFiled] = time.Now().UnixMilli()
	encoded := EncodeParams(params)
	var fullUrl = fmt.Sprintf("%s%s?%s", RestAPI, route, encoded)

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
	var response = new(bytes.Buffer)

	if _, err := io.Copy(response, resp.Body); err != nil {
		return err
	}

	fmt.Printf("response(%d) %s\n", resp.StatusCode, response.String())
	//errCode := utils.Json.Get(response.Bytes(), "code").ToInt()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		err := utils.Json.Unmarshal(response.Bytes(), &errResp)
		if err != nil {
			return err
		}
		return fmt.Errorf("binance request error(%d)[%s]", errResp.Code, errResp.Msg)
	}

	return utils.Json.Unmarshal(response.Bytes(), &returnType)
}
