package binance

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

func (c *Connector) GetCandles(symbol, interval string, limit int64) (types.CandlesEntry, error) {
	var klines [][]any
	err := c.Call(http.MethodGet, KlineEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"interval":  interval,
		"limit":     limit,
	}, constants.None, &klines)
	if err != nil {
		return nil, err
	}
	var candles = make(types.CandlesEntry, len(klines))
	for i, kline := range klines {
		list := append(kline[:6], kline[7])
		candles[i] = make([]float64, len(list))
		for j, v := range list {
			candles[i][j] = types.Safe2Float(v)
		}
	}
	return candles, nil
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, resp)
	return resp.ServerTime, err
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (types.OrderBookEntry, error) {
	var limit int64 = 30
	if depth != nil {
		limit = *depth
	}
	var orderBook = new(struct {
		LastUpdateId int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	})
	err := c.Call(http.MethodGet, DepthEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"limit":     limit,
	}, constants.None, orderBook)
	if err != nil {
		return types.OrderBookEntry{}, err
	}
	return types.OrderBookEntry{
		Symbol:    symbol,
		Bids:      orderBook.Bids,
		Asks:      orderBook.Asks,
		Timestamp: orderBook.LastUpdateId,
	}, nil
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp = new(struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	})
	err := c.Call(http.MethodGet, PriceTickerEndpoint, &platforms.ObjectBody{
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
