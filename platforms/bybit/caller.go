package bybit

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
	"strconv"
	"time"
)

func (c *Connector) Sign(params []byte) string {
	mac := hmac.New(sha256.New, []byte(c.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Connector) Call(method string, route string, params platforms.Serializer, authType constants.AuthType, returnType any) error {
	// Add necessary parameters
	var body io.Reader
	bodyData, err := params.Serialize()
	if err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	recvWindow := "5000"
	var signData = new(bytes.Buffer)
	signData.WriteString(timestamp)
	signData.WriteString(c.APIKey)
	signData.WriteString(recvWindow)
	var header = http.Header{}
	header.Add("User-Agent", "cex.connector/1.5")

	url := RestAPI + route
	if method == http.MethodPost {
		//signData.Write(bodyData)
		_, err = signData.ReadFrom(bodyData)
		if err != nil {
			return err
		}
		body = bodyData
	} else {
		queryString, err := params.EncodeQuery()
		if err != nil {
			return err
		}
		signData.WriteString(queryString)
		url = fmt.Sprintf("%s?%s", url, queryString)
	}
	if authType == constants.Signed {
		// must sign
		header.Set(signTypeKey, "2")
		header.Set(apiRequestKey, c.APIKey)
		header.Set(timestampKey, timestamp)
		header.Set(recvWindowKey, recvWindow)
		signature := c.Sign(signData.Bytes())
		header.Set(signatureKey, signature)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header = header
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("[Bitmart] Response %s", resp.Status)
	}
	return json.Unmarshal(respBody, returnType)
}
