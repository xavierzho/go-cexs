package bitmart

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
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
	err := c.Call(http.MethodGet, BalanceEndpoint, map[string]interface{}{}, constants.Keyed, &response)
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

func (c *Connector) GetOrderBook(symbol string, limit *int64) (*types.OrderBookEntry, error) {
	var response OrderBookResponse
	if limit == nil {
		*limit = 30
	}
	err := c.Call(http.MethodGet, OrderBookEndpoint, map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}, constants.None, &response)
	if err != nil {
		return nil, err
	}
	ts, err := strconv.ParseInt(response.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	return &types.OrderBookEntry{
		Symbol:    symbol,
		Asks:      response.Asks,
		Bids:      response.Bids,
		Timestamp: ts,
	}, nil
}

func (c *Connector) GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error) {
	var resp [][]any
	err := c.Call(http.MethodGet, KlineEndpoint, map[string]interface{}{
		"symbol": symbol,
		"step":   interval,
		"limit":  limit,
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume", "volume_usd",
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
	err := c.Call(http.MethodGet, ServerTimeEndpoint, map[string]interface{}{}, constants.None, resp)
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
	err := c.Call(http.MethodGet, TickerEndpoint, map[string]interface{}{
		"symbol": symbol,
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
