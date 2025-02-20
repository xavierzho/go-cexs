package mexc

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"time"
)

func (c *Connector) GetOrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error) {
	var resp = new(struct {
		LastUpdateId int        `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	})
	if depth == nil {
		*depth = 30
	}
	err := c.Call(http.MethodGet, OrderBookEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"limit":     depth,
	}, constants.None, resp)
	if err != nil {
		return nil, err
	}
	return &types.OrderBookEntry{
		Bids:      resp.Bids,
		Asks:      resp.Asks,
		Symbol:    symbol,
		Timestamp: time.Now().Unix(),
	}, nil
}

func (c *Connector) GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error) {
	var resp [][]any
	err := c.Call(http.MethodGet, CandleEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"interval":  interval,
		"limit":     limit,
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}

	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume", "time_close", "value_usd",
	}
	var result []types.CandleEntry
	for _, kline := range resp {
		var candle = new(types.CandleEntry)

		candle.FromList(kline, keys)
		result = append(result, *candle)
	}
	return result, nil
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64 `json:"server_time"`
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, resp)
	if err != nil {
		return 0, err
	}
	return resp.ServerTime, nil
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp = new(types.TickerEntry)

	err := c.Call(http.MethodGet, TickerEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.None, resp)
	return *resp, err
}
