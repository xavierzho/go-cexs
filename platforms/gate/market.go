package gate

import (
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

type OrderBook struct {
	ID      string     `json:"id"`
	Current string     `json:"current"`
	Update  string     `json:"update"`
	Asks    [][]string `json:"asks"`
	Bids    [][]string `json:"bids"`
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (types.OrderBookEntry, error) {
	if depth == nil {
		*depth = 30
	}
	var resp OrderBook
	err := c.Call(http.MethodGet, QueryOrderBookEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"limit":     depth,
	}, constants.None, &resp)
	if err != nil {
		return types.OrderBookEntry{}, err
	}
	timestamp, _ := strconv.ParseInt(resp.Update, 10, 64)
	return types.OrderBookEntry{
		Asks:      resp.Asks,
		Bids:      resp.Bids,
		Timestamp: timestamp,
	}, nil
}

func (c *Connector) GetCandles(symbol, interval string, limit int64) (types.CandlesEntry, error) {
	var resp [][]any
	err := c.Call(http.MethodGet, QueryCandleEndpoint, &platforms.ObjectBody{
		"currency_pair": c.SymbolPattern(symbol),
		"interval":      interval,
		"limit":         limit,
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	var result = make(types.CandlesEntry, len(resp))
	for i, k := range resp {
		var line = []any{k[0], k[5], k[4], k[3], k[2], k[6], k[1]}
		result[i] = make(types.CandleEntry, len(line))
		for j, a := range line {
			v := types.Safe2Float(a)
			if j == 0 {
				v *= 1000
			}
			result[i][j] = v
		}
	}
	return result, err
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp = new(struct {
		ServerTime int64 `json:"server_time"`
	})
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, resp)
	return resp.ServerTime, err
}

type Ticker struct {
	ChangeUtc8       string `json:"change_utc8"`
	Last             string `json:"last"`
	QuoteVolume      string `json:"quote_volume"`
	BaseVolume       string `json:"base_volume"`
	EtfLeverage      string `json:"etf_leverage"`
	EtfNetValue      string `json:"etf_net_value"`
	HighestBid       string `json:"highest_bid"`
	EtfPreNetValue   string `json:"etf_pre_net_value"`
	CurrencyPair     string `json:"currency_pair"`
	ChangePercentage string `json:"change_percentage"`
	ChangeUtc0       string `json:"change_utc0"`
	High24h          string `json:"high_24h"`
	Low24h           string `json:"low_24h"`
	EtfPreTimestamp  int    `json:"etf_pre_timestamp"`
	LowestAsk        string `json:"lowest_ask"`
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp Ticker
	err := c.Call(http.MethodGet, QueryTickerEndpoint, &platforms.ObjectBody{
		"timezone":  "utc8",
		SymbolFiled: symbol,
	}, constants.None, &resp)
	if err != nil {
		return types.TickerEntry{}, nil
	}
	price, _ := decimal.NewFromString(resp.Last)

	return types.TickerEntry{
		Symbol: resp.CurrencyPair,
		Price:  price,
	}, nil
}
