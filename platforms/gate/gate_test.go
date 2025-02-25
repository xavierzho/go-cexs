package gate

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"testing"
	"time"
)

func TestSign(t *testing.T) {

	var signBytes = new(bytes.Buffer)
	signBytes.WriteString("GET")
	signBytes.WriteByte('\n')
	signBytes.WriteString("/api/v4/futures/orders")
	signBytes.WriteByte('\n')
	signBytes.WriteString("contract=BTC_USD&status=finished&limit=50")
	signBytes.WriteByte('\n')
	signBytes.WriteString("cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e")
	signBytes.WriteByte('\n')
	signBytes.WriteString("1541993715")

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
