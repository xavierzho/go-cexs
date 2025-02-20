package mexc

import (
	"context"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"regexp"
	"strings"
)

type MarketStream struct {
	*platforms.StreamBase
}

type Event interface {
	GetSymbol() string
}
type StreamResp[T Event] struct {
	Symbol    string `json:"s"`
	Timestamp string `json:"t"`
	Operator  string `json:"c"`
	Data      T      `json:"d"`
}

type DepthUpdate struct {
	Bids    []Depth `json:"bids"`
	Asks    []Depth `json:"asks"`
	Event   string  `json:"e"`
	Version string  `json:"r"`
}

func (d DepthUpdate) GetSymbol() string {
	return d.Event
}

type Depth struct {
	Price  string `json:"p"`
	Volume string `json:"v"`
}

func (stream *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error {
	err := stream.Connect(StreamAPI)
	if err != nil {
		return err
	}
	err = stream.SendMessage(map[string]any{
		"method": SubscribeOp,
		"params": []string{
			fmt.Sprintf("spot@public.limit.depth.v3.api@%s@20", strings.ToUpper(symbol)),
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
				var resp StreamResp[DepthUpdate]
				_ = utils.Json.Unmarshal(msg, &resp)
				var asks = make([][]string, len(resp.Data.Asks))
				var bids = make([][]string, len(resp.Data.Bids))
				for i, bid := range resp.Data.Bids {
					bids[i] = []string{bid.Price, bid.Volume}
				}
				for i, ask := range resp.Data.Asks {
					asks[i] = []string{ask.Price, ask.Volume}
				}
				channel <- types.DepthEntry{
					Asks: asks,
					Bids: bids,
				}
			}
		}
	}()
	return nil
}

type CandleUpdate struct {
	Kline Candle `json:"k"`
	Event string `json:"e"`
}

type Candle struct {
	Volume    string `json:"a"`
	Close     string `json:"c"`
	TimeEnd   int64  `json:"T"`
	TimeStart int64  `json:"t"`
	Quantity  string `json:"v"`
	High      string `json:"h"`
	Interval  string `json:"i"`
	Low       string `json:"l"`
	Open      string `json:"o"`
}

func (d CandleUpdate) GetSymbol() string {
	return d.Event
}
func itl(i string) (string, error) {
	re := regexp.MustCompile("(\\d+)([mshdwM])")
	match := re.FindStringSubmatch(i)
	if len(match) != 3 {
		return "", fmt.Errorf("invalid time symbol")
	}
	value := match[1]
	unit := match[2]
	switch unit {
	case "m":
		return fmt.Sprintf("Min%s", value), nil
	case "h":
		return fmt.Sprintf("Hour%s", value), nil
	case "d":
		return fmt.Sprintf("Day%s", value), nil
	case "w":
		return fmt.Sprintf("Week%s", value), nil
	case "M":
		return fmt.Sprintf("Mounth%s", value), nil
	default:
		return "", fmt.Errorf("not support %s unit", unit)
	}
}
func (stream *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error {
	err := stream.Connect(StreamAPI)
	if err != nil {
		return err
	}
	interval, err = itl(interval)
	if err != nil {
		return err
	}
	err = stream.SendMessage(map[string]any{
		"method": SubscribeOp,
		"params": []string{
			fmt.Sprintf("spot@public.kline.v3.api@%s@%s", symbol, interval),
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
				var resp StreamResp[CandleUpdate]
				_ = utils.Json.Unmarshal(msg, &resp)
				var kline = resp.Data.Kline
				channel <- types.CandleEntry{
					"open":       kline.Open,
					"high":       kline.High,
					"low":        kline.Low,
					"close":      kline.Close,
					"time_start": kline.TimeStart,
					"time_end":   kline.TimeEnd,
					"volume":     kline.Volume,
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
