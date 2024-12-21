package bybit

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client

	Channels map[string]*websocket.Conn
}

func (c *Connector) SymbolPattern(symbol string) string {
	symbol, _ = constants.StandardizeSymbol(symbol)
	return symbol
}

func (c *Connector) GetOrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) GetCandles(symbol, interval string, limit int64) ([]types.CandleEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) GetServerTime() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) GetTicker(symbol string) (types.TickerEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Sign(body []byte) string {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Call(method string, route string, body map[string]interface{}, authType constants.AuthType, returnType interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) BatchOrder(orders []types.OrderEntry) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {
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

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) OrderBook(symbol string, depth *int64) (*types.OrderBookEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Connector) Name() constants.Platform {
	return constants.ByBit
}

func NewConnector(base *platforms.Credentials, client *http.Client) platforms.Connector {
	return &Connector{Credentials: base, Client: client}
}
