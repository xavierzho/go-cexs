package types

import "github.com/xavierzho/go-cexs/constants"

type OrderStatusConverter interface {
	String() string
	Convert() constants.OrderStatus
}

type OrderTypeConverter interface {
	String() string
	Convert() constants.OrderType
}
