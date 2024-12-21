package binance

import (
	"github.com/shopspring/decimal"
	"net/http"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

func (c *Connector) GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error) {
	var klines [][]any
	err := c.Call(http.MethodGet, KlineEndpoint, map[string]any{
		SymbolFiled: symbol,
		"interval":  interval,
		"limit":     limit,
	}, constants.None, &klines)
	if err != nil {
		return nil, err
	}
	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume", "time_end",
		"volume_usd", "trades", "base_buy", "quote_buy", "ignore",
	}
	var candles []types.CandleEntry
	for _, kline := range klines {
		var candle = make(types.CandleEntry)

		candle.FromList(kline, keys)
		candles = append(candles, candle)
	}
	return candles, nil
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, map[string]any{}, constants.None, resp)
	return resp.ServerTime, err
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error) {
	var limit int64 = 30
	if depth != nil {
		limit = *depth
	}
	var orderBook = new(struct {
		LastUpdateId int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	})
	err := c.Call(http.MethodGet, DepthEndpoint, map[string]any{
		SymbolFiled: symbol,
		"limit":     limit,
	}, constants.None, orderBook)

	return &types.OrderBookEntry{
		Symbol:    symbol,
		Bids:      orderBook.Bids,
		Asks:      orderBook.Asks,
		Timestamp: orderBook.LastUpdateId,
	}, err
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp = new(struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	})
	err := c.Call(http.MethodGet, PriceTickerEndpoint, map[string]any{
		SymbolFiled: symbol,
	}, constants.None, resp)
	if err != nil {
		return types.TickerEntry{}, err
	}
	price, _ := decimal.NewFromString(resp.Price)
	return types.TickerEntry{
		Symbol: symbol,
		Price:  price,
	}, nil
}
