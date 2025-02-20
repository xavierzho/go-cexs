package binance

import "github.com/xavierzho/go-cexs/constants"

const StreamAPI = "wss://stream.binance.com:9443/stream"

const RestAPI = "https://api.binance.com"

type OrderStatus string

func (s OrderStatus) String() string {
	return string(s)
}
func (s OrderStatus) Convert() constants.OrderStatus {
	switch s {
	case OrderStatusFilled:
		return constants.Filled
	case OrderStatusCanceled:
		return constants.Canceled
	case OrderStatusExpired:
		return constants.Canceled
	case OrderStatusNew:
		return constants.Open
	case OrderStatusPartiallyFilled:
		return constants.PartiallyFilled
	case OrderStatusPendingCanceled:
		return constants.PartiallyCanceled
	case OrderStatusPendingNew:
		return constants.Open
	case OrderStatusRejected:
		return constants.Canceled
	case OrderStatusExpiredInMatch:
		return constants.Canceled
	default:
		return constants.Error
	}
}

const (
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPendingNew      OrderStatus = "PENDING_NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusPendingCanceled OrderStatus = "PENDING_CANCEL"
	OrderStatusRejected        OrderStatus = "REJECTED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
	OrderStatusExpiredInMatch  OrderStatus = "EXPIRED_IN_MATCH"
)

type OrderType string

const (
	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
)

func (ot OrderType) String() string {
	return string(ot)
}
func (ot OrderType) Convert() constants.OrderType {
	switch ot {
	case OrderTypeLimit:
		return constants.Limit
	case OrderTypeMarket:
		return constants.Market
	case OrderTypeLimitMaker:
		return constants.LimitMaker
	case OrderTypeStopLoss:
		return constants.StopLoss
	case OrderTypeStopLossLimit:
		return constants.StopLossLimit
	case OrderTypeTakeProfit:
		return constants.TakeProfit
	case OrderTypeTakeProfitLimit:
		return constants.TakeProfitLimit
	default:
		return constants.Market
	}
}

const (
	OrderEndpoint       = "/api/v3/order"
	OpenOrdersEndpoint  = "/api/v3/openOrders"
	DepthEndpoint       = "/api/v3/depth"
	AccountEndpoint     = "/api/v3/account"
	ServerTimeEndpoint  = "/api/v3/time"
	KlineEndpoint       = "/api/v3/klines"
	PriceTickerEndpoint = "/api/v3/ticker/price"
	ListenKeyEndpoint   = "/api/v3/userDataStream"
)

type NewOrderRespType string

func (t NewOrderRespType) String() string {
	return string(t)
}

const (
	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeResult NewOrderRespType = "RESULT"
	NewOrderRespTypeFULL   NewOrderRespType = "FULL"
)

const (
	SymbolFiled    = "symbol"
	TimeFiled      = "timestamp"
	SignatureFiled = "signature"
	HeaderAPIKEY   = "X-MBX-APIKEY"
)

type TimeInForce string

const (
	GTC TimeInForce = "GTC"
)

func (t TimeInForce) String() string {
	return string(t)
}

type EventType string

const (
	OrderEventType   EventType = "executionReport"
	AccountEventType EventType = "outboundAccountPosition"
	BalanceEventType EventType = "balanceUpdate"
	ExpiredEventType EventType = "listenKeyExpired"
)
