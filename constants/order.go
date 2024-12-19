package constants

type OrderType int

const (
	Limit OrderType = iota
	Market
	LimitMaker
	StopLoss
	StopLossLimit
	TakeProfit
	TakeProfitLimit
)

type OrderStatus int

const (
	Canceled = -1
	Open     = iota
	PartiallyFilled
	Filled
	PartiallyCanceled
	Error = -2
)
