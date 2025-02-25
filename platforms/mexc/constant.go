package mexc

import "github.com/xavierzho/go-cexs/constants"

const StreamAPI = "wss://wbs.mexc.com/ws"

const RestAPI = "https://api.mexc.com"

const (
	KeyHeader = "X-MEXC-APIKEY"

	SymbolFiled = "symbol"

	SubscribeOp = "SUBSCRIPTION"

	UnSubscribeOP = "UNSUBSCRIPTION"
)

const (
	ServerTimeEndpoint = "/api/v3/time"
	OrderBookEndpoint  = "/api/v3/depth"
	CandleEndpoint     = "/api/v3/klines"
	TickerEndpoint     = "/api/v3/ticker/price"
	OrderEndpoint      = "/api/v3/order"
	BatchOrderEndpoint = "/api/v3/batchOrders"
	OpenOrdersEndpoint = "/api/v3/openOrders"
	AccountEndpoint    = "/api/v3/account"
	ListenKeyEndpoint  = "/api/v3/userDataStream"
)

type OrderType string

const (
	OrderTypeLimit       OrderType = "LIMIT"
	OrderTypeMarket      OrderType = "MARKET"
	OrderTypeLimitMarker OrderType = "LIMIT_MAKER"
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
	case OrderTypeLimitMarker:
		return constants.LimitMaker
	default:
		return constants.Market
	}
}

type OrderStatus string

const (
	OrderStatusNew               OrderStatus = "NEW"
	OrderStatusFilled            OrderStatus = "FILLED"
	OrderStatusPartiallyFilled   OrderStatus = "PARTIALLY_FILLED"
	OrderStatusCanceled          OrderStatus = "CANCELED"
	OrderStatusPartiallyCanceled OrderStatus = "PARTIALLY_CANCELED"
)

func (o OrderStatus) String() string {
	return string(o)
}

func (o OrderStatus) Convert() constants.OrderStatus {
	switch o {
	case OrderStatusNew:
		return constants.Open
	case OrderStatusFilled:
		return constants.Filled
	case OrderStatusPartiallyFilled:
		return constants.PartiallyFilled
	case OrderStatusCanceled:
		return constants.Canceled
	case OrderStatusPartiallyCanceled:
		return constants.PartiallyCanceled
	default:
		return constants.Error
	}
}

type ChangeType string

const (
	Deposit          ChangeType = "DEPOSIT"
	DepositFee       ChangeType = "DEPOSIT_FEE"
	ContractTransfer ChangeType = "CONTRACT_TRANSFER"
	InternalTransfer ChangeType = "INTERNAL_TRANSFER"
	Withdraw         ChangeType = "WITHDRAW"
	WithdrawFee      ChangeType = "WITHDRAW_FEE"
	Entrust          ChangeType = "ENTRUST"
	EntrustPlace     ChangeType = "ENTRUST_PLACE"
	EntrustCancel    ChangeType = "ENTRUST_CANCEL"
	TradeFee         ChangeType = "TRADE_FEE"
	EntrustUnfrozen  ChangeType = "ENTRUST_UNFROZEN"
	Airdrop          ChangeType = "SUGAR"
	EtfIndex         ChangeType = "ETF_INDEX"
)
