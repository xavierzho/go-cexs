package okx

import "github.com/xavierzho/go-cexs/constants"

const RestAPI = "https://www.okx.com"

const StreamAPI = "wss://ws.okx.com:8443"

const (
	PublicChannel   = "/ws/v5/public"
	PrivateChannel  = "/ws/v5/private"
	BusinessChannel = "/ws/v5/business"

	OrderEndpoint               = "/api/v5/trade/order"
	OrderBatchEndpoint          = "/api/v5/trade/batch-orders"
	ServerTimeEndpoint          = "/api/v5/public/time"
	CandleRealTimeEndpoint      = "/api/v5/market/candles"
	CandleHistoryEndpoint       = "/api/v5/market/history-index-candles"
	TickerEndpoint              = "/api/v5/market/index-tickers"
	OrderBookEndpoint           = "/api/v5/market/books"
	OrderCancelEndpoint         = "/api/v5/trade/cancel-order"
	OrderCancelBatchEndpoint    = "/api/v5/trade/cancel-batch-orders"
	OrderPendingEndpoint        = "/api/v5/trade/orders-pending"
	OrderCancelAllAfterEndpoint = "/api/v5/trade/cancel-all-after"
	AccountBalanceEndpoint      = "/api/v5/account/balance"
)

type TradeMode string

const (
	CashMode         TradeMode = "cash"
	SpotIsolatedMode TradeMode = "spot_isolated"
	CrossMode        TradeMode = "cross"
	IsolatedMode     TradeMode = "isolated"
)

type OrderType string

const (
	OrderTypeLimit       OrderType = "limit"
	OrderTypeMarket      OrderType = "market"
	OrderTypeOnlyMaker   OrderType = "post_only"
	OrderTypeFok         OrderType = "fok"
	OrderTypeIoc         OrderType = "ioc"
	OrderTypeLimitIoc    OrderType = "optimal_limit_ioc" // Market order with immediate-or-cancel order (applicable only to Expiry Futures and Perpetual Futures).
	OrderTypeMMP         OrderType = "mmp"               // Market Maker Protection (only applicable to Option in Portfolio Margin mode)
	OrderTypeMMPPostOnly OrderType = "mmp_and_post_only" // Market Maker Protection and Post-only order(only applicable to Option in Portfolio Margin mode)
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
	case OrderTypeOnlyMaker:
		return constants.LimitMaker
	default:
		return constants.Limit
	}
}

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "live"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusMMpCanceled     OrderStatus = "mmp_canceled"
)

func (o OrderStatus) String() string {
	return string(o)
}

func (o OrderStatus) Convert() constants.OrderStatus {
	switch o {
	case OrderStatusNew:
		return constants.Open
	case OrderStatusCanceled, OrderStatusMMpCanceled:
		return constants.Canceled
	case OrderStatusPartiallyFilled:
		return constants.PartiallyFilled
	case OrderStatusFilled:
		return constants.Filled
	default:
		return constants.Error
	}
}
