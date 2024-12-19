package bybit

import (
	"github.com/gorilla/websocket"
	cexconns "github.com/xavierzho/go-cexs"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
)

type Connector struct {
	*cexconns.Credentials
	Client *http.Client

	Channels map[string]*websocket.Conn
}

func (c *Connector) Sign(r *http.Request, body []byte) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Call(method string, route string, body map[string]interface{}, authType constants.AuthType, returnType interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) PlaceOrder(params types.UnifiedOrder) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) BatchOrder(orders []types.UnifiedOrder) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) GetOrderStatus(symbol string, orderId string) (*constants.OrderStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) CancelAll(symbol string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Balance(symbols []string) (map[string]types.UnifiedBalance, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) OrderBook(symbol string, depth *int64) (*types.UnifiedOrderBook, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) PendingOrders(symbol string) ([]types.UnifiedOpenOrder, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Name() constants.Platform {
	return constants.ByBit
}

func NewConnector(base *cexconns.Credentials, client *http.Client) cexconns.Connector {
	return &Connector{Credentials: base, Client: client}
}
