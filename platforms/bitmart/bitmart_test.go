package bitmart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"log"
	"testing"
	"time"
)

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

func TestMarketStream(t *testing.T) {
	stream := NewMarketStream()
	ctx, cancel := context.WithCancel(context.Background())
	var symbol = "BTCUSDT"
	var candles = make(chan types.CandleEntry)
	err := stream.CandleStream(ctx, symbol, "1m", candles)
	if err != nil {
		t.Error(err)
	}
	var depths = make(chan types.DepthEntry)
	stream2 := NewMarketStream()
	err = stream2.DepthStream(ctx, symbol, depths)
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case candle := <-candles:
				fmt.Println(symbol, candle)
			case depth := <-depths:
				fmt.Println(symbol, depth)
			}
		}
	}()
	time.Sleep(10 * time.Second)
	cancel()
}

func TestGetCandles(t *testing.T) {

	var symbol = "BTCUSDT"
	opt := ""
	cex := NewConnector(&platforms.Credentials{
		APIKey:    "",
		APISecret: "",
		Option:    &opt,
	}, nil)
	candles, err := cex.GetCandles(symbol, "1m", 200)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(candles)
}
