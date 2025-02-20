package types

import "github.com/xavierzho/go-cexs/constants"

// OrderStatusConverter Convert the status of the corresponding exchange to a unified status
type OrderStatusConverter interface {
	String() string
	Convert() constants.OrderStatus
}

// OrderTypeConverter Convert the order type of the corresponding exchange to a unified order type
type OrderTypeConverter interface {
	String() string
	Convert() constants.OrderType
}
