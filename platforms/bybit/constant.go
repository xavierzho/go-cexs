package bybit

import (
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
)

const StreamAPI = "wss://stream.bybit.com"
const StreamAPITestnet = "wss://stream-testnet.bybit.com"
const RestAPI = "https://api.bybit.com"

const (
	// Globals
	timestampKey  = "X-BAPI-TIMESTAMP"
	signatureKey  = "X-BAPI-SIGN"
	apiRequestKey = "X-BAPI-API-KEY"
	recvWindowKey = "X-BAPI-RECV-WINDOW"
	signTypeKey   = "X-BAPI-SIGN-TYPE"
)

const (
	ServerTimeEndpoint       = "/v5/market/time"
	CandleEndpoint           = "/v5/market/kline"
	OrderBookEndpoint        = "/v5/market/orderbook"
	TickerEndpoint           = "/v5/market/tickers"
	PlaceOrderEndpoint       = "/v5/order/create"
	BatchPlaceOrderEndpoint  = "/v5/order/create-batch"
	RealTimeOrderEndpoint    = "/v5/order/realtime"
	OrderCancelEndpoint      = "/v5/order/cancel"
	OrderCancelAllEndpoint   = "/v5/order/cancel-all"
	OrderBatchCancelEndpoint = "/v5/order/cancel-batch"
	WalletBalanceEndpoint    = "/v5/account/wallet-balance"
)
const (
	SpotMainnetChannel            = "/v5/public/spot"
	PerpetualMainnetChannel       = "/v5/public/linear"
	USDCMainnetChannel            = "v5/public/option"
	InverseContractMainnetChannel = "/v5/public/inverse"
	PrivateChannel                = "/v5/private"
)

type RestResp[T, Ext fmt.Stringer] struct {
	Code   int    `json:"retCode"`
	Msg    string `json:"retMsg"`
	Result T      `json:"result"`
	Ext    Ext    `json:"retExtInfo"`
	Time   int64  `json:"time"`
}

type OrderType string

const (
	OrderTypeLimit  OrderType = "Limit"
	OrderTypeMarket OrderType = "Market"
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
	default:
		return constants.Market
	}
}

type OrderStatus string

const (
	OrderStatusOpen              OrderStatus = "New"
	OrderStatusFilled            OrderStatus = "Filled"
	OrderStatusPartiallyFilled   OrderStatus = "PartiallyFilled"
	OrderStatusCanceled          OrderStatus = "Cancelled"
	OrderStatusPartiallyCanceled OrderStatus = "PartiallyFilledCanceled"
	OrderStatusRejected          OrderStatus = "Rejected"
	OrderStatusTriggered         OrderStatus = "Triggered"
	OrderStatusUnTriggered       OrderStatus = "Untriggered"
	OrderStatusDeactivated       OrderStatus = "Deactivated"
)

func (o OrderStatus) String() string {
	return string(o)
}
func (o OrderStatus) Convert() constants.OrderStatus {
	switch o {
	case OrderStatusOpen, OrderStatusTriggered:
		return constants.Open
	case OrderStatusFilled:
		return constants.Filled
	case OrderStatusPartiallyFilled:
		return constants.PartiallyFilled
	case OrderStatusCanceled, OrderStatusDeactivated, OrderStatusRejected, OrderStatusUnTriggered:
		return constants.Canceled
	case OrderStatusPartiallyCanceled:
		return constants.PartiallyCanceled
	default:
		return constants.Error
	}
}

type AccountType string

const (
	UnifiedAccount  AccountType = "UNIFIED"
	SpotAccount     AccountType = "SPOT"
	ContractAccount AccountType = "CONTRACT"
)
