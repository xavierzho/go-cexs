package okx

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"strconv"
)

type OrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp string     `json:"ts"`
}

func (o OrderBook) String() string {
	return o.Timestamp
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (types.OrderBookEntry, error) {
	var resp RestReturn[OrderBook]
	if depth == nil {
		*depth = 30
	}
	err := c.Call(http.MethodGet, OrderBookEndpoint, &platforms.ObjectBody{
		"instId": symbol,
		"sz":     *depth,
	}, constants.None, &resp)
	if err != nil {
		return types.OrderBookEntry{}, err
	}
	ts, _ := strconv.ParseInt(resp.Data[0].Timestamp, 10, 64)
	return types.OrderBookEntry{
		Asks:      resp.Data[0].Asks,
		Bids:      resp.Data[0].Bids,
		Symbol:    symbol,
		Timestamp: ts,
	}, nil
}

type Candle []any

func (Candle) String() string {
	return ""
}
func (c *Connector) GetCandles(symbol, interval string, limit int64) (types.CandlesEntry, error) {
	var resp RestReturn[Candle]

	err := c.Call(http.MethodGet, CandleRealTimeEndpoint, &platforms.ObjectBody{
		"instId": symbol,
		"limit":  limit,
		"bar":    interval,
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	var result = make(types.CandlesEntry, len(resp.Data))
	for i, k := range resp.Data {
		list := append(k[:7])
		result[i] = make(types.CandleEntry, len(list))
		for j, a := range list {
			result[i][j] = types.Safe2Float(a)
		}
	}
	return result, nil
}

type ServerTime struct {
	Timestamp string `json:"ts"`
}

func (s ServerTime) String() string {
	return s.Timestamp
}
func (c *Connector) GetServerTime() (int64, error) {
	var resp RestReturn[ServerTime]
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, &resp)
	if err != nil {
		return 0, err
	}

	ts, err := strconv.ParseInt(resp.Data[0].Timestamp, 10, 64)
	if err != nil {
		return 0, err
	}
	return ts, nil
}

type Ticker struct {
	Symbol    string `json:"instId"`
	Price     string `json:"idxPx"`
	High24h   string `json:"high24h"`
	SodUtc0   string `json:"sodUtc0"`
	Open24h   string `json:"open24h"`
	Low24h    string `json:"low24h"`
	SodUtc8   string `json:"sodUtc8"`
	Timestamp string `json:"ts"`
}

func (t Ticker) String() string {
	return t.Symbol
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp RestReturn[Ticker]
	err := c.Call(http.MethodGet, TickerEndpoint, &platforms.ObjectBody{
		"instId": symbol,
	}, constants.None, &resp)
	if err != nil {
		return types.TickerEntry{}, err
	}
	price, _ := decimal.NewFromString(resp.Data[0].Price)
	return types.TickerEntry{
		Symbol: symbol,
		Price:  price,
	}, nil
}
