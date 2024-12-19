package okx

import (
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
)

func TestWsConnect(t *testing.T) {
	const url = "wss://ws.okx.com:8443/ws/v5/public"
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}
	conn.WriteJSON(map[string]any{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"channel": "index-tickers",
				"instId":  "BTC-USDT",
			},
		},
	})
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("subscibe success %s\n", msg)
	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				t.Error(err)
			}

			fmt.Printf("recv %s\n", msg)
		}
	}()
	select {}
}
