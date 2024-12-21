package bitmart

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/platforms"
	"io/ioutil"
	"log"
	"testing"
)

func zipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
func TestWsConnect(t *testing.T) {
	const url = "wss://ws-manager-compress.bitmart.com/api?protocol=1.1"
	const private = "wss://ws-manager-compress.bitmart.com/user?protocol=1.1"
	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}
	var subPayload = map[string]any{
		"op": "subscribe",
		"args": []string{
			"spot/ticker:DOGE_USDT",
		},
	}

	payload, err := json.Marshal(subPayload)
	if err != nil {
		t.Error(err)
	}

	conn.WriteMessage(websocket.TextMessage, payload)
	_, subResp, err := conn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	done := make(chan struct{})
	fmt.Printf("subscribe resp %s \n", subResp)
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

func TestStreamSign(t *testing.T) {
	timestamp := "1589267764859"
	apikey := "80618e45710812162b04892c7ee5ead4a3cc3e56"
	apisecert := "6c6c98544461bbe71db2bca4c6d7fd0021e0ba9efc215f9c6ad41852df9d9df9"
	apimemo := "test001"
	signature := "3ceeb7e1b8cb165a975e28a2e2dfaca4d30b358873c0351c1a071d8c83314556"
	stream := NewUserStream(platforms.NewCredentials(apikey, apisecert, &apimemo))
	signed := stream.Sign(timestamp)
	if signed != signature {
		t.Errorf("no match %s-%s", signed, signature)
	}
}
