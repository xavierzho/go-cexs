package types

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
)

type UnifiedOrder struct {
	Symbol      string              `json:"symbol"`
	Type        constants.OrderType `json:"type"`
	Side        string              `json:"side"`
	Price       decimal.Decimal     `json:"price"`
	Quantity    decimal.Decimal     `json:"quantity"`
	TradeNo     string              `json:"trade_no"`
	TimeInForce *string             `json:"-"`
}

type UnifiedBalance struct {
	Free     string `json:"free"`
	Locked   string `json:"locked"`
	Currency string `json:"currency"`
}

type UnifiedOrderBook struct {
	Symbol    string     `json:"symbol"`
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp int64      `json:"timestamp"`
}
type UnifiedOpenOrder struct {
	Symbol   string                `json:"symbol"`
	Type     constants.OrderType   `json:"type"`
	Side     string                `json:"side"`
	Price    decimal.Decimal       `json:"price"`
	Quantity decimal.Decimal       `json:"quantity"`
	TradeNo  string                `json:"trade_no"`
	Status   constants.OrderStatus `json:"status"`
	OrderId  string                `json:"order_id"`
}
