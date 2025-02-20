package gate

import "github.com/xavierzho/go-cexs/constants"

const StreamAPI = "wss://api.gateio.ws/ws/v4/"

const RestAPI = "https://api.gateio.ws"

const (
	SignHeader      = "SIGN"
	KeyHeader       = "KEY"
	TimestampHeader = "Timestamp"

	SymbolFiled = "trading_pair"
)

const (
	APIPrefix              = "/api/v4"
	BatchOrdersEndpoint    = APIPrefix + "/spot/batch_orders"
	OrderEndpoint          = APIPrefix + "/spot/orders"
	QueryTickerEndpoint    = APIPrefix + "/spot/tickers"
	QueryOrderBookEndpoint = APIPrefix + "/spot/order_book"
	QueryCandleEndpoint    = APIPrefix + "/spot/candlesticks"
	ServerTimeEndpoint     = APIPrefix + "/spot/time"
	BatchCancelEndpoint    = APIPrefix + "/spot/cancel_batch_orders"
	OpenOrdersEndpoint     = APIPrefix + "/spot/open_orders"
	SmallBalanceEndpoint   = APIPrefix + "/wallet/small_balance"
)

type OrderStatus string

const (
	OrderStatusOpen     OrderStatus = "open"
	OrderStatusClosed   OrderStatus = "closed"
	OrderStatusCanceled OrderStatus = "canceled"
)

func (o OrderStatus) String() string {
	return string(o)
}
func (o OrderStatus) Convert() constants.OrderStatus {
	switch o {
	case OrderStatusOpen:
		return constants.Open
	case OrderStatusClosed:
		return constants.Filled
	case OrderStatusCanceled:
		return constants.Canceled
	default:
		return constants.Error
	}
}

type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

func (o OrderType) String() string {
	return string(o)
}

func (o OrderType) Convert() constants.OrderType {
	switch o {
	case OrderTypeLimit:
		return constants.Limit
	case OrderTypeMarket:
		return constants.Market
	default:
		return constants.Market
	}
}

const (
	TimeInForceGTC = "gtc"
	TimeInForceIOC = "ioc"
	TimeInForcePOC = "poc"
	TimeInForceFOK = "fok"
)
