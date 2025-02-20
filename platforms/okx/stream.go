package okx

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
type StreamEvent[T fmt.Stringer] struct {
	Data []T `json:"data"`
	Arg  struct {
		InstID  string `json:"instId"`
		Channel string `json:"channel"`
	} `json:"arg,omitempty"`
	Action string `json:"action"`
	Code   string `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	ConnId string `json:"connId"`
}
type DepthEvent struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp string     `json:"ts"`
	Checksum  int64      `json:"checksum"`
	PrevSeqId int64      `json:"prevSeqId"`
	SeqId     int64      `json:"seqId"`
}

func (e DepthEvent) String() string {
	return ""
}
func (stream *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error {
	err := stream.Connect(fmt.Sprintf("%s%s", StreamAPI, PublicChannel))
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithHyphen(symbol)

	err = stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []map[string]any{
			{"channel": "books", "instId": symbol},
		},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = stream.Close()
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamEvent[DepthEvent]

				_ = utils.Json.Unmarshal(msg, &event)
				for _, e := range event.Data {
					channel <- types.DepthEntry{
						Asks: e.Asks,
						Bids: e.Bids,
					}
				}
			}
		}
	}()
	return nil
}

type CandleEvent []any

func (CandleEvent) String() string {
	return ""
}

func (stream *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error {
	err := stream.Connect(fmt.Sprintf("%s%s", StreamAPI, BusinessChannel))
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithHyphen(symbol)
	err = stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []map[string]any{
			{"channel": "candle" + interval, "instId": symbol},
		},
	})
	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume", "volume_usd", "volume_usd", "is_close",
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = stream.Close()
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamEvent[CandleEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				for _, e := range event.Data {
					var kline = new(types.CandleEntry)
					kline.FromList(e, keys)
					channel <- *kline
				}
			}
		}
	}()
	return nil
}

func NewMarketStream() platforms.MarketStreamer {
	return &MarketStream{
		StreamBase: platforms.NewStream(),
	}
}
