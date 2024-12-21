package platforms

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
)

type Auth interface {
	Sign(r *http.Request, body []byte)
	Call(method string, route string, body map[string]interface{},
		authType constants.AuthType, returnType interface{}) error
}

type Caller interface {
	Sign(params []byte) string
	Call(method string, route string, body map[string]interface{},
		authType constants.AuthType, returnType interface{}) error
}
type Connector interface {
	Caller
	Name() constants.Platform
	SymbolPattern(symbol string) string
	Trade
	MarketData
}

type MarketData interface {
	GetOrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error)
	GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error)
	GetServerTime() (int64, error)
	GetTicker(symbol string) (types.TickerEntry, error)
}
type Trade interface {
	PlaceOrder(params types.OrderEntry) (string, error)
	BatchOrder(orders []types.OrderEntry) ([]string, error)
	GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error)

	Cancel(symbol, orderId string) (bool, error)
	CancelAll(symbol string) error
	CancelByIds(symbol string, orderIds []string) (map[string]bool, error)

	Balance(symbols []string) (map[string]types.BalanceEntry, error)
	PendingOrders(symbol string) ([]types.OpenOrderEntry, error)
}

type Credentials struct {
	APIKey    string
	APISecret string
	Option    *string
}

func NewCredentials(apikey, apiSecret string, option *string) *Credentials {
	return &Credentials{
		APIKey:    apikey,
		APISecret: apiSecret,
		Option:    option,
	}
}
