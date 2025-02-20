package constants

// OrderType unified order type
type OrderType int

const (
	Limit OrderType = iota
	Market
	LimitMaker
	StopLoss
	StopLossLimit
	TakeProfit
	TakeProfitLimit
	Iceberg
)

// OrderStatus unified order status
type OrderStatus int

const (
	Canceled = -1
	Open     = iota
	PartiallyFilled
	Filled
	PartiallyCanceled
	Error = -2
)
