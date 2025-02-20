package mexc

import (
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"testing"
	"time"

	"github.com/xavierzho/go-cexs/types"
)

func TestMarketAPI(t *testing.T) {
	cex := NewConnector(&platforms.Credentials{}, &http.Client{})
	var symbol = "BTCUSDT"
	candles, err := cex.GetCandles(symbol, "1m", 300)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(candles)
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

func TestUserDataStream(t *testing.T) {

	stream := NewUserStream(platforms.NewCredentials("mx0vglqaFSgIoT4AG1", "69b258384be14fa3a6401370eb1c94d5", nil))
	stream.Login()
}
