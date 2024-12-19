package binance

import (
	"net/http"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

func (c *Connector) Candles(symbol, interval string, limit int64) ([]types.Candle, error) {
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
	var candles []types.Candle
	for _, kline := range klines {
		var candle = make(types.Candle)

		candle.FromList(kline, keys)
		candles = append(candles, candle)
	}
	return candles, nil
}

func (c *Connector) ServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, map[string]any{}, constants.None, resp)
	return resp.ServerTime, err
}

func (c *Connector) OrderBook(symbol string, depth *int64) (*types.UnifiedOrderBook, error) {
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

	return &types.UnifiedOrderBook{
		Symbol:    symbol,
		Bids:      orderBook.Bids,
		Asks:      orderBook.Asks,
		Timestamp: orderBook.LastUpdateId,
	}, err
}
