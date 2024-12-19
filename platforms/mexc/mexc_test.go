package mexc

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func TestWsConnect(t *testing.T) {
	const baseURL = "wss://wbs.mexc.com/ws"
	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(baseURL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	var subPayload = map[string]any{
		"method": "SUBSCRIPTION", "params": []string{
			"spot@public.deals.v3.api@BTCUSDT",
		},
	}
	conn.WriteJSON(subPayload)
	_, pong, err := conn.ReadMessage()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("sub resp %s\n", pong)
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
