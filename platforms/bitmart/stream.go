package bitmart

import (
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
)

type MarketStream struct {
	*platforms.StreamBase
}

func NewMarketStream() *MarketStream {
	return &MarketStream{
		StreamBase: platforms.NewStream(),
	}
}

type Event interface {
	GetSymbol() string
}

type StreamResp[T Event] struct {
	Table string `json:"table"`
	Data  []T    `json:"data"`
}

type CandleUpdate struct {
	Candle []any  `json:"candle"`
	Symbol string `json:"symbol"`
}

func (c CandleUpdate) GetSymbol() string {
	return c.Symbol
}

func (stream *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error {
	err := stream.Connect(PublicChannel)
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithUnderline(symbol)
	err = stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []string{
			fmt.Sprintf("spot/kline%s:%s", interval, symbol),
		},
	})
	if err != nil {
		return err
	}

	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume",
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResp[CandleUpdate]

				_ = utils.Json.Unmarshal(msg, &event)
				for _, datum := range event.Data {
					var candle = new(types.CandleEntry)
					candle.FromList(datum.Candle, keys)
					channel <- *candle
				}

			}
		}
	}()
	return nil
}

type Depth struct {
	Symbol  string     `json:"symbol"`
	MsT     int64      `json:"ms_t"`
	Type    string     `json:"type"`
	Version int        `json:"version"`
	Asks    [][]string `json:"asks"`
	Bids    [][]string `json:"bids"`
}

func (d Depth) GetSymbol() string {
	return d.Symbol
}
func (stream *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error {
	err := stream.Connect(PublicChannel)
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithUnderline(symbol)
	err = stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []string{
			fmt.Sprintf("spot/depth/increase100:%s", symbol),
		},
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResp[Depth]

				_ = utils.Json.Unmarshal(msg, &event)
				channel <- types.DepthEntry{
					Asks: event.Data[0].Asks,
					Bids: event.Data[0].Bids,
				}
			}
		}
	}()
	return nil
}
