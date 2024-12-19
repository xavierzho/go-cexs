package binance

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"strings"
)

type MarketStream struct {
	dialer *websocket.Dialer
	conn   *websocket.Conn
}

type CandleEvent struct {
	Symbol string `json:"s"`
	StreamEvent
	Kline Kline `json:"k"`
}

type Kline struct {
	B            string `json:"B,omitempty"`
	Close        string `json:"c"`
	FirstOrderId int64  `json:"f"`
	High         string `json:"h"`
	Interval     string `json:"i"`
	LastOrderId  int64  `json:"L"`
	Low          string `json:"l"`
	NumOfTrades  int    `json:"n"`
	Open         string `json:"o"`
	Qty          string `json:"q"`
	QuoteVolume  string `json:"Q"` // Taker buy quote asset volume
	Symbol       string `json:"s"`
	StartTime    int64  `json:"t"`
	EndTime      int64  `json:"T"`
	Volume       string `json:"v"`
	BaseVolume   string `json:"V"` // Taker buy base asset volume
	IsClose      bool   `json:"x"`
}

type DepthEvent struct {
	StreamEvent
	Symbol  string     `json:"s"`
	FirstId int64      `json:"u"`
	LastId  int64      `json:"U"`
	Bids    [][]string `json:"b"`
	Asks    [][]string `json:"a"`
}

func (stream *MarketStream) Connect(ctx context.Context) error {

	conn, _, err := stream.dialer.DialContext(ctx, StreamAPI, nil)
	if err != nil {
		return err
	}

	stream.conn = conn
	return nil
}

func (stream *MarketStream) CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.Candle) {
	err := stream.conn.WriteJSON(map[string]any{
		"method": "SUBSCRIBE",
		"id":     uuid.New().String(),
		"params": []string{
			fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval),
		},
	})
	if err != nil {
		return
	}

	if _, _, err := stream.conn.ReadMessage(); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := stream.conn.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResponse[CandleEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				k := event.Data.Kline
				channel <- types.Candle{
					"open":       k.Open,
					"high":       k.High,
					"low":        k.Low,
					"close":      k.Close,
					"volume":     k.Volume,
					"time_start": k.StartTime,
					"time_end":   k.EndTime,
				}
			}
		}
	}()
}

func (stream *MarketStream) DepthStream(ctx context.Context, symbol string, channel chan<- DepthEvent) {
	err := stream.conn.WriteJSON(map[string]any{
		"method": "SUBSCRIBE",
		"id":     uuid.New().String(),
		"params": []string{
			fmt.Sprintf("%s@depth", strings.ToLower(symbol)),
		},
	})

	if err != nil {
		return
	}

	if _, _, err = stream.conn.ReadMessage(); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := stream.conn.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResponse[DepthEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				channel <- event.Data
			}
		}
	}()
}