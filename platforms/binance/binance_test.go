package binance

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"testing"

	"github.com/xavierzho/go-cexs/platforms"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

func TestWsConnect(t *testing.T) {
	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial("wss://stream.binance.com:9443/stream", nil)
	if err != nil {
		t.Error(err)
	}
	//var subPayload = map[string]any{
	//	"method": "SUBSCRIBE",
	//	"id":     "1",
	//	"params": []string{
	//		"pepeusdt@ticker",
	//	},
	//}
	//
	//conn.WriteJSON(subPayload)
	//_, subResp, err := conn.ReadMessage()
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Printf("subscribe resp %s \n", subResp)
	conn.WriteJSON(map[string]any{
		"method": "SUBSCRIBE",
		"id":     "2",
		"params": []string{
			"btcusdt@kline_1m",
		},
	})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			buf := bytes.NewBuffer(message)

			log.Printf("recv: %s", buf.String())
		}
	}()
	select {}
}

func TestWsLogin(t *testing.T) {
	apikey := os.Getenv("BinanceAPIKEY")
	secret := os.Getenv("BinanceSERCET")

	stream := NewUserStream(platforms.NewCredentials(apikey, secret, nil))
	err := stream.Login()
	if err != nil {
		t.Error(err)
	}

	_ = stream.Login()

}
func encodeValues(v url.Values) []byte {
	if v == nil {
		return nil
	}
	var buf []byte
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		if len(vs) == 0 {
			continue
		}
		if len(buf) > 0 {
			buf = append(buf, '&')
		}
		buf = append(buf, url.QueryEscape(k)...)
		buf = append(buf, '=')
		if len(vs) == 1 {
			buf = append(buf, url.QueryEscape(vs[0])...)
			continue
		}
		vss, _ := json.Marshal(&vs)
		buf = append(buf, url.QueryEscape(string(vss))...)
	}
	return buf
}
func TestSign(t *testing.T) {
	//var apiKey = "vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A"
	var secretKey = "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"
	query, _ := url.ParseQuery("symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559")
	//query := SortParams(map[string]any{
	//	"symbol":      "BTCUSDT",
	//	"side":        "SELL",
	//	"type":        "LIMIT",
	//	"timeInForce": "GTC",
	//	"quantity":    "1.0000000",
	//	"price":       "0.20",
	//})
	var qs = encodeValues(query)

	var queryString = "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559"
	if bytes.NewBuffer(qs).String() != queryString {
		t.Errorf("not match query string between (%s, %s)", bytes.NewBuffer(qs).String(), queryString)
		return
	}
	h := hmac.New(sha256.New, bytes.NewBufferString(secretKey).Bytes())
	h.Write(bytes.NewBufferString(queryString).Bytes())

	signature := hex.EncodeToString(h.Sum(nil))

	if signature != "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71" {
		t.Errorf("not match %s - c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71", signature)
	}
}

func TestOrder(t *testing.T) {
	cred := &platforms.Credentials{
		APIKey:    os.Getenv("BinanceAPIKEY"),
		APISecret: os.Getenv("BinanceSERCET"),
	}
	fmt.Println(cred)
	connector := NewConnector(cred, nil)
	var symbol = "VTHOUSDT"
	orderId, err := connector.PlaceOrder(types.OrderEntry{
		Symbol:   symbol,
		Type:     constants.Limit,
		Side:     "SELL",
		Quantity: decimal.NewFromFloat(2463),
		Price:    decimal.NewFromFloat(0.0032),
		TradeNo:  uuid.New().String(),
	})
	if err != nil {
		t.Error(err)
	}

	status, err := connector.GetOrderStatus(symbol, orderId)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(orderId, status)
	_, err = connector.Cancel(symbol, orderId)
	if err != nil {
		t.Error(err)
	}
}
