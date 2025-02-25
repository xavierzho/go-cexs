package bitmart

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"net/http"
	"strconv"
)

type BalanceResponse struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var response struct {
		Wallet []BalanceResponse
	}
	err := c.Call(http.MethodGet, BalanceEndpoint, &platforms.ObjectBody{}, constants.Keyed, &response)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]types.BalanceEntry, len(symbols))
	for _, balance := range response.Wallet {
		result[balance.ID] = types.BalanceEntry{
			Currency: balance.ID,
			Free:     balance.Available,
			Locked:   balance.Frozen,
		}
	}
	return result, nil
}

type OrderBookResponse struct {
	Timestamp string     `json:"ts"`     // Create time(Timestamp in milliseconds)
	Symbol    string     `json:"symbol"` // Trading pair
	Asks      [][]string `json:"asks"`   // Order book on sell side
	Bids      [][]string `json:"bids"`   // Order book on buy side
	Amount    string     `json:"amount"` // Total number of current price depth
	Price     string     `json:"price"`  // The price at current depth
}

func (c *Connector) GetOrderBook(symbol string, limit *int64) (types.OrderBookEntry, error) {
	var response OrderBookResponse
	if limit == nil {
		*limit = 30
	}
	err := c.Call(http.MethodGet, OrderBookEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"limit":     limit,
	}, constants.None, &response)
	if err != nil {
		return types.OrderBookEntry{}, err
	}
	ts, err := strconv.ParseInt(response.Timestamp, 10, 64)
	if err != nil {
		return types.OrderBookEntry{}, err
	}
	return types.OrderBookEntry{
		Symbol:    symbol,
		Asks:      response.Asks,
		Bids:      response.Bids,
		Timestamp: ts,
	}, nil
}

func (c *Connector) GetCandles(symbol, interval string, limit int64) (types.CandlesEntry, error) {
	var klines [][]any
	sec, err := utils.ToSeconds(interval)
	if err != nil {
		return nil, err
	}
	err = c.Call(http.MethodGet, KlineEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"step":      sec / 60,
		"limit":     limit,
	}, constants.None, &klines)
	if err != nil {
		return nil, err
	}

	var candles = make(types.CandlesEntry, len(klines))
	for i, kline := range klines {
		candles[i] = make([]float64, len(kline))
		for j, v := range kline {
			candles[i][j] = types.Safe2Float(v)
		}
	}
	return candles, nil
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64 `json:"server_time"`
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, resp)
	return resp.ServerTime, err
}

type TickerResp struct {
	Symbol      string `json:"symbol"`
	AskSz       string `json:"ask_sz"`
	AskPx       string `json:"ask_px"`
	Last        string `json:"last"`
	Qv24h       string `json:"qv_24h"`
	V24h        string `json:"v_24h"`
	High24h     string `json:"high_24h"`
	Low24h      string `json:"low_24h"`
	BidSz       string `json:"bid_sz"`
	BidPx       string `json:"bid_px"`
	Fluctuation string `json:"fluctuation"`
	Open24h     string `json:"open_24h"`
	Ts          string `json:"ts"`
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp = new(TickerResp)
	err := c.Call(http.MethodGet, TickerEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.None, resp)
	if err != nil {
		return types.TickerEntry{}, err
	}
	price, _ := decimal.NewFromString(resp.Last)
	return types.TickerEntry{
		Symbol: symbol,
		Price:  price,
	}, nil
}
