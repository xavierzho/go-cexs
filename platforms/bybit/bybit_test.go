package bybit

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/utils"
	"log"
	"testing"
)

const apiKey = "3jvDbQruYPwRQ2BL8d"
const apiSecret = "q77TJ1Z6LptjcOVkc8g2dF9Z4ZJE63OTT1tL"

func TestConnectStream(t *testing.T) {
	var conns = websocket.DefaultDialer
	conn, _, err := conns.Dial("wss://stream.bybit.com/v5/public/spot", nil)
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})

	var topic = map[string]any{
		"req_id": "test", // 可選
		"op":     "subscribe",
		"args": []string{
			"tickers.PEPEUSDT",
		},
	}

	topicMessage, err := utils.Json.Marshal(topic)
	if err != nil {
		panic(err)
	}
	var depthTopic = map[string]any{
		"req_id": "",
		"op":     "subscribe",
		"args": []string{
			"orderbook.1.PEPEUSDT",
		},
	}
	conn.WriteMessage(websocket.TextMessage, topicMessage)
	_, msg, _ := conn.ReadMessage()
	log.Printf("%+v, %s \n", topic, msg)

	depthMessage, err := utils.Json.Marshal(depthTopic)
	if err != nil {
		panic(err)
	}
	conn.WriteMessage(websocket.TextMessage, depthMessage)
	_, msg, _ = conn.ReadMessage()
	log.Printf("%+v, %s \n", depthTopic, msg)
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
