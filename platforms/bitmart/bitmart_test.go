package bitmart

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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
