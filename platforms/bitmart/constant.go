package bitmart

import "github.com/xavierzho/go-cexs/constants"

const PublicChannel = "wss://ws-manager-compress.bitmart.com/api?protocol=1.1"

const PrivateChannel = "wss://ws-manager-compress.bitmart.com/user?protocol=1.1"

const RestAPI = "https://api-cloud.bitmart.com"

type OrderType string

const (
	OrderTypeMarket     OrderType = "market"
	OrderTypeLimit      OrderType = "limit"
	OrderTypeLimitMaker OrderType = "limit_maker_market"
)

func (o OrderType) String() string {
	return string(o)
}

func (o OrderType) Convert() constants.OrderType {
	switch o {
	case OrderTypeMarket:
		return constants.Market
	case OrderTypeLimit:
		return constants.Limit
	case OrderTypeLimitMaker:
		return constants.LimitMaker
	default:
		return constants.Market
	}
}

type OrderStatus string

const (
	OrderStatusNew               = "new"
	OrderStatusCanceled          = "canceled"
	OrderStatusFilled            = "filled"
	OrderStatusPartiallyFilled   = "partially_filled"
	OrderStatusPartiallyCanceled = "partially_canceled"
)

func (o OrderStatus) String() string {
	return string(o)
}

func (o OrderStatus) Convert() constants.OrderStatus {
	switch o {
	case OrderStatusNew:
		return constants.Open
	case OrderStatusPartiallyFilled:
		return constants.PartiallyFilled
	case OrderStatusFilled:
		return constants.Filled
	case OrderStatusCanceled:
		return constants.Canceled
	case OrderStatusPartiallyCanceled:
		return constants.PartiallyCanceled
	default:
		return constants.Error
	}
}

const (
	CancelAllEndpoint  = "/spot/v4/cancel_all"
	CancelEndpoint     = "/spot/v3/cancel_order"
	CancelsEndpoint    = "/spot/v4/cancel_orders"
	BalanceEndpoint    = "/spot/v1/wallet"
	OpenOrdersEndpoint = "/spot/v4/query/open-orders"
	OrderBookEndpoint  = "/spot/quotation/v3/books"
	NewOrderEndpoint   = "/spot/v2/submit_order"
	BatchOrderEndpoint = "/spot/v4/batch_orders"
	QueryOrderEndpoint = "/spot/v4/query/order"
	TickerEndpoint     = "/spot/quotation/v3/ticker"
	ServerTimeEndpoint = "/system/time"
	KlineEndpoint      = "/spot/quotation/v3/lite-klines"
)

const (
	SymbolFiled = "symbol"
)
