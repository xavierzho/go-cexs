package bybit

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"strconv"
)

type NullExt struct {
}

func (NullExt) String() string {
	return ""
}

type OrderBook struct {
	Symbol    string     `json:"s"`
	Asks      [][]string `json:"a"`
	Bids      [][]string `json:"b"`
	Timestamp int64      `json:"ts"`
	Seq       int64      `json:"seq"`
	Cts       int64      `json:"cts"`
	U         int64      `json:"u"`
}

func (OrderBook) String() string {
	return ""
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error) {
	var resp RestResp[OrderBook, NullExt]
	if depth == nil {
		*depth = 30
	}
	err := c.Call(http.MethodGet, OrderBookEndpoint, &platforms.ObjectBody{
		"category": "spot",
		"symbol":   symbol,
		"limit":    *depth,
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	return &types.OrderBookEntry{
		Symbol:    resp.Result.Symbol,
		Asks:      resp.Result.Asks,
		Bids:      resp.Result.Bids,
		Timestamp: resp.Result.Timestamp,
	}, nil
}

type Candle struct {
	Symbol   string  `json:"symbol"`
	Category string  `json:"category"`
	List     [][]any `json:"list"`
}

func (Candle) String() string {
	return ""
}
func (c *Connector) GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error) {
	var ret RestResp[Candle, NullExt]
	err := c.Call(http.MethodGet, CandleEndpoint, &platforms.ObjectBody{
		"category": "spot",
		"symbol":   symbol,
		"interval": timeConvert(interval),
		"limit":    limit,
	}, constants.None, &ret)
	if err != nil {
		return nil, err
	}
	var keys = []string{
		"time_start", "open", "high", "low", "close", "volume", "turnover",
	}
	var result []types.CandleEntry
	for _, kline := range ret.Result.List {
		var candle = new(types.CandleEntry)
		candle.FromList(kline, keys)
		result = append(result, *candle)
	}
	return result, nil
}

type ServerTime struct {
	Second string `json:"timeSecond"`
	Nano   string `json:"timeNano"`
}

func (ServerTime) String() string {
	return ""
}

func (c *Connector) GetServerTime() (int64, error) {
	var resp RestResp[ServerTime, NullExt]
	err := c.Call(http.MethodGet, ServerTimeEndpoint, &platforms.ObjectBody{}, constants.None, &resp)
	if err != nil {
		return 0, err
	}
	t, err := strconv.ParseInt(resp.Result.Second, 10, 64)
	if err != nil {
		return 0, err
	}
	return t, nil
}

type Ticker struct {
	Symbol                 string `json:"symbol"`
	Bid1Price              string `json:"bid1Price"`
	IndexPrice             string `json:"indexPrice"`
	OpenInterest           string `json:"openInterest"`
	DeliveryTime           string `json:"deliveryTime"`
	LowPrice24h            string `json:"lowPrice24h"`
	OpenInterestValue      string `json:"openInterestValue"`
	BasisRate              string `json:"basisRate"`
	Volume24h              string `json:"volume24h"`
	NextFundingTime        string `json:"nextFundingTime"`
	Turnover24h            string `json:"turnover24h"`
	PredictedDeliveryPrice string `json:"predictedDeliveryPrice"`
	Bid1Size               string `json:"bid1Size"`
	Basis                  string `json:"basis"`
	MarkPrice              string `json:"markPrice"`
	PrevPrice1h            string `json:"prevPrice1h"`
	PrevPrice24h           string `json:"prevPrice24h"`
	Ask1Size               string `json:"ask1Size"`
	Price24hPcnt           string `json:"price24hPcnt"`
	HighPrice24h           string `json:"highPrice24h"`
	DeliveryFeeRate        string `json:"deliveryFeeRate"`
	Ask1Price              string `json:"ask1Price"`
	LastPrice              string `json:"lastPrice"`
	FundingRate            string `json:"fundingRate"`
}

type Tickers struct {
	Category string   `json:"category"`
	List     []Ticker `json:"list"`
}

func (Tickers) String() string {
	return ""
}
func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	var resp RestResp[Tickers, NullExt]
	err := c.Call(http.MethodGet, TickerEndpoint, &platforms.ObjectBody{"symbol": symbol}, constants.None, &resp)
	if err != nil {
		return types.TickerEntry{}, err
	}
	price, _ := decimal.NewFromString(resp.Result.List[0].LastPrice)
	return types.TickerEntry{
		Symbol: symbol,
		Price:  price,
	}, nil
}
