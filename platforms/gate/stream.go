package gate

import (
	"context"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"time"
)

type MarketStream struct {
	*platforms.StreamBase
}
type Event[T any] struct {
	Result  T      `json:"result"`
	TimeMs  int64  `json:"time_ms"`
	Channel string `json:"channel"`
	Time    int    `json:"time"`
	Event   string `json:"event"`
}

type DepthUpdate struct {
	S           string     `json:"s"`
	Timestamp   int64      `json:"t"`
	EventType   string     `json:"e"`
	EventTime   int        `json:"E"`
	FirstUpdate int        `json:"U"`
	LastUpdate  int        `json:"u"`
	Bids        [][]string `json:"b"`
	Asks        [][]string `json:"a"`
}

func (m *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error {
	err := m.Connect(StreamAPI)
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithUnderline(symbol)
	err = m.SendMessage(map[string]any{
		"time":    time.Now().Unix(),
		"channel": "spot.order_book_update",
		"event":   "subscribe",
		"payload": []string{
			symbol, "100ms",
		},
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				m.Close()
				return
			default:
				msg, err := m.ReadMessage()
				if err != nil {
					continue
				}
				var event Event[DepthUpdate]
				err = utils.Json.Unmarshal(msg, &event)
				if err != nil {
					continue
				}

				channel <- types.DepthEntry{
					Asks: event.Result.Asks,
					Bids: event.Result.Bids,
				}
			}
		}
	}()
	return nil
}

type CandleUpdate struct {
	A string `json:"a"`
	C string `json:"c"`
	T string `json:"t"`
	V string `json:"v"`
	W bool   `json:"w"`
	H string `json:"h"`
	L string `json:"l"`
	N string `json:"n"`
	O string `json:"o"`
}

func (m *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error {
	err := m.Connect(StreamAPI)
	if err != nil {
		return err
	}
	symbol = constants.SymbolWithUnderline(symbol)
	err = m.SendMessage(map[string]any{
		"time":    time.Now().Unix(),
		"channel": "spot.candlesticks",
		"event":   "subscribe",
		"payload": []string{
			interval, symbol,
		},
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				m.Close()
				return
			default:
				msg, err := m.ReadMessage()
				if err != nil {
					continue
				}
				var event Event[CandleUpdate]
				err = utils.Json.Unmarshal(msg, &event)
				if err != nil {
					continue
				}
				channel <- types.CandleEntry{
					"open_time": event.Result.T,
					"open":      event.Result.O,
					"high":      event.Result.H,
					"low":       event.Result.L,
					"close":     event.Result.C,
					"volume":    event.Result.V,
					"is_close":  event.Result.W,
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
