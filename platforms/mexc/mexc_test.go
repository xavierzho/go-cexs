package mexc

import (
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"os"
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
	apikey := os.Getenv("MexcAPIKEY")
	secret := os.Getenv("MexcSECRET")
	//fmt.Println(apikey, secret)
	stream := NewUserStream(platforms.NewCredentials(apikey, secret, nil))
	err := stream.Login()
	if err != nil {

		t.Errorf("login failed: %s", err)
		return
	}
	stream1 := NewUserStream(platforms.NewCredentials(apikey, secret, nil))
	err = stream1.Login()
	if err != nil {
		t.Errorf("login1 failed: %s", err)
	}
	//var balanceChan = make(chan types.BalanceUpdateEntry)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//err = stream.BalanceStream(ctx, balanceChan)
	//if err != nil {
	//	fmt.Println("balance stream", err)
	//	return
	//}
	//for {
	//	select {
	//	case bal := <-balanceChan:
	//		fmt.Println(bal)
	//
	//	}
	//}
	var orderChan = make(chan types.OrderUpdateEntry)
	err = stream.OrderStream(ctx, orderChan)
	if err != nil {
		fmt.Println("order stream", err)
		return
	}
	var placeChan = make(chan types.OrderUpdateEntry)
	err = stream1.(*UserDataStream).PlaceStream(ctx, placeChan)
	if err != nil {
		fmt.Println("place stream", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-orderChan:
			fmt.Println("order", order)
			continue
		case place := <-placeChan:
			fmt.Println("place", place)
			continue
			//case bal := <-balanceChan:
			//	fmt.Println(bal)

		}
	}
}
