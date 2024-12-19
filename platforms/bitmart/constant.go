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