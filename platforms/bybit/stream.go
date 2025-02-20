package bybit

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
)

type MarketStream struct {
	*platforms.StreamBase
}

type PublicStream[T fmt.Stringer] struct {
	Topic     string `json:"topic"`
	Type      string `json:"type"`
	Timestamp int64  `json:"ts"`
	Data      T      `json:"data"`
}
type DepthEvent struct {
	Symbol   string     `json:"s"`
	Bids     [][]string `json:"b"`
	Asks     [][]string `json:"a"`
	Seq      int64      `json:"seq"`
	UpdateId int64      `json:"u"`
}

func (e DepthEvent) String() string {
	return e.Symbol
}

func (m *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error {
	err := m.Connect(fmt.Sprintf("%s%s", StreamAPI, SpotMainnetChannel))
	if err != nil {
		return err
	}

	err = m.SendMessage(map[string]any{
		"op":     "subscribe",
		"req_id": uuid.New().String(),
		"args": []string{
			fmt.Sprintf("orderbook.%d.%s", 200, symbol),
		},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = m.Close()
				return
			default:
				msg, err := m.ReadMessage()
				if err != nil {
					continue
				}
				var event PublicStream[DepthEvent]
				err = utils.Json.Unmarshal(msg, &event)
				if err != nil {
					continue
				}
				channel <- types.DepthEntry{
					Bids: event.Data.Bids,
					Asks: event.Data.Asks,
				}
			}
		}
	}()
	return nil
}

type CandleEvent []struct {
	Volume    string `json:"volume"`
	Confirm   bool   `json:"confirm"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Interval  string `json:"interval"`
	Close     string `json:"close"`
	Turnover  string `json:"turnover"`
	Open      string `json:"open"`
	Timestamp int64  `json:"timestamp"`
}

func (e CandleEvent) String() string {
	return ""
}

func (m *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error {
	err := m.Connect(fmt.Sprintf("%s%s", StreamAPI, SpotMainnetChannel))
	if err != nil {
		return err
	}
	err = m.SendMessage(map[string]any{
		"op":     "subscribe",
		"req_id": uuid.New().String(),
		"args": []string{
			fmt.Sprintf("kline.%s.%s", timeConvert(interval), symbol),
		},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = m.Close()
				return
			default:
				msg, err := m.ReadMessage()
				if err != nil {
					continue
				}
				var event PublicStream[CandleEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				for _, e := range event.Data {
					channel <- types.CandleEntry{
						"open":       e.Open,
						"high":       e.High,
						"low":        e.Low,
						"close":      e.Close,
						"time_start": e.Start,
						"time_end":   e.End,
						"volume":     e.Volume,
					}

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
