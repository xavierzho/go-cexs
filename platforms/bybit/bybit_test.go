package bybit

import (
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"testing"
	"time"

	"github.com/xavierzho/go-cexs/types"
)

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
	cred := &platforms.Credentials{
		APIKey:    "",
		APISecret: "",
	}
	cex := NewConnector(cred, http.DefaultClient)

	candles, err := cex.GetCandles("BTCUSDT", "1m", 200)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(candles)
}
