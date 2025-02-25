package okx

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

func (c *Connector) Call(method string, route string, params platforms.Serializer, _ constants.AuthType, returnType interface{}) error {
	// Add necessary parameters
	var body io.Reader
	var err error
	symbol, ok := params.Exists("instId")
	if ok {
		params.Set("instId", c.SymbolPattern(symbol.(string)))
	}
	headers := http.Header{}
	headers.Set("OK-ACCESS-KEY", c.APIKey)
	headers.Set("OK-ACCESS-PASSPHRASE", *c.Option)
	timestamp := time.Now().Format("2006-01-02T15:04:05.999Z")
	headers.Set("OK-ACCESS-TIMESTAMP", timestamp)
	headers.Set("Content-Type", "application/json")
	prevSign := fmt.Sprintf("%s%s%s", timestamp, method, route)
	url := RestAPI + route
	if method == http.MethodGet {
		query, err := params.EncodeQuery()
		if err != nil {
			return err
		}
		q := fmt.Sprintf("?%s", query)
		prevSign += q
		url += q
	} else if method == http.MethodPost {
		var bodyBytes = new(bytes.Buffer)
		body, err = params.Serialize()
		if err != nil {
			return err
		}

		_, err = io.Copy(bodyBytes, body)
		if err != nil {
			return err
		}

		prevSign += fmt.Sprintf("%s", bodyBytes.Bytes())
	}
	headers.Set("OK-ACCESS-SIGN", c.Sign([]byte(prevSign)))
	fmt.Println(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header = headers
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Okx] bad status code: %d", resp.StatusCode)
	}
	return json.Unmarshal(respBody, returnType)
}
