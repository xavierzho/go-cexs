package platforms

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

// Caller defines the interface for making authenticated API calls.
type Caller interface {
	// Sign generates a signature for the given parameters.
	Sign(params []byte) string
	// Call makes an API call to the specified endpoint.
	// method: HTTP method (e.g., GET, POST).
	// route: API endpoint path.
	// body: Request body (for POST requests).
	// authType: Authentication type required for the call.
	// returnType: Pointer to the struct where the response data will be unmarshalled.
	Call(method string, route string, body Serializer,
		authType constants.AuthType, returnType interface{}) error
}

// SpotConnector defines the interface for interacting with a specific exchange.
type SpotConnector interface {
	// Caller Embeds interface for authentication.
	Caller
	// Name returns the name of the exchange.
	Name() constants.Platform
	// SymbolPattern formats a symbol into a standardized pattern (e.g., BTC_USDT).
	SymbolPattern(symbol string) string
	// Trade Embeds interface for trading operations.
	Trade
	// SpotMarketData Embeds interface for market data retrieval.
	SpotMarketData
}

// SpotMarketData defines the interface for retrieving market data.
type SpotMarketData interface {
	// GetOrderBook retrieves the order book for a given symbol.
	// symbol: Trading pair symbol (e.g., BTCUSDT).
	// depth: Number of order book entries to retrieve (optional).
	GetOrderBook(symbol string, depth *int64) (types.OrderBookEntry, error)
	// GetCandles retrieves candlestick data (OHLCV) for a given symbol and interval.
	// symbol: Trading pair symbol (e.g., BTCUSDT).
	// interval: Time interval for the candles (e.g., 1m, 5m, 1h, 1d).
	// limit: Maximum number of candles to retrieve.
	GetCandles(symbol, interval string, limit int64) (types.CandlesEntry, error)
	// GetServerTime retrieves the server time of the exchange.
	GetServerTime() (int64, error)
	// GetTicker retrieves the ticker information for a given symbol.
	// symbol: Trading pair symbol (e.g., BTCUSDT).
	GetTicker(symbol string) (types.TickerEntry, error)
}

// Trade defines the interface for trading operations.
type Trade interface {
	// PlaceOrder places a new order.
	// params: Order parameters.
	PlaceOrder(params types.OrderEntry) (string, error)
	// BatchOrder places multiple orders at once.
	// orders: A slice of order parameters.
	BatchOrder(orders []types.OrderEntry) ([]string, error)
	// QueryOrder retrieve order
	// symbol: Trading pair symbol.
	// orderId: ID of the order.
	QueryOrder(symbol string, orderId string) (types.QueryOrder, error)
	// GetOrderStatus just retrieve the status of a specific order.
	// symbol: Trading pair symbol.
	// orderId: ID of the order.
	GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error)
	// Cancel cancels a specific order.
	// symbol: Trading pair symbol.
	// orderId: ID of the order.
	Cancel(symbol, orderId string) (bool, error)
	// CancelAll cancels all pending orders for a given symbol.
	// symbol: Trading pair symbol.
	CancelAll(symbol string) error
	// CancelByIds cancels orders by their IDs.
	// symbol: Trading pair symbol.
	// orderIds: A slice of order IDs.
	CancelByIds(symbol string, orderIds []string) (map[string]bool, error)
	// Balance retrieves account balances. If symbols is empty, retrieves all balances.
	// symbols: A slice of symbols to retrieve balances for (optional).
	Balance(symbols []string) (map[string]types.BalanceEntry, error)
	// PendingOrders retrieves all pending (open) orders for a given symbol.
	// symbol: Trading pair symbol.
	PendingOrders(symbol string) ([]types.OpenOrderEntry, error)
}

// Credentials each exchange oauth keys
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
