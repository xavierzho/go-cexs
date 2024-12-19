package cexconns

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
	Auth
	Name() constants.Platform

	SpotConnector
}

type MarketConnector interface {
	OrderBook(symbol string, depth *int64) (*types.UnifiedOrderBook, error)
	Candles(symbol, interval string, limit int64) ([]types.Candle, error)
	ServerTime() (int64, error)
}
type SpotConnector interface {
	PlaceOrder(params types.UnifiedOrder) (string, error)
	BatchOrder(orders []types.UnifiedOrder) ([]string, error)
	GetOrderStatus(symbol string, orderId string) (*constants.OrderStatus, error)

	Cancel(symbol, orderId string) (bool, error)
	CancelAll(symbol string) error
	CancelByIds(symbol string, orderIds []string) (map[string]bool, error)

	Balance(symbols []string) (map[string]types.UnifiedBalance, error)
	OrderBook(symbol string, depth *int64) (*types.UnifiedOrderBook, error)
	PendingOrders(symbol string) ([]types.UnifiedOpenOrder, error)
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

func Exchange(ex string, apikey, apiSecret string, option *string) Connector {
	switch constants.Platform(ex) {
	//case constants.Bitmart:
	//return bitmart.NewConnector(NewCredentials(apikey, apiSecret, option), &http.Client{})
	default:
		return nil
	}
}
