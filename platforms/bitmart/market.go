package bitmart

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"strconv"
	"strings"
)

type BalanceResponse struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

func (c *Connector) Balance(symbols []string) (map[string]types.UnifiedBalance, error) {
	var response struct {
		Wallet []BalanceResponse
	}
	err := c.Call(http.MethodGet, "spot/v1/wallet", map[string]interface{}{}, constants.Keyed, &response)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]types.UnifiedBalance, len(symbols))
	for _, balance := range response.Wallet {
		result[balance.ID] = types.UnifiedBalance{
			Currency: balance.ID,
			Free:     balance.Available,
			Locked:   balance.Frozen,
		}
	}
	return result, nil
}

func (c *Connector) PendingOrders(symbol string) ([]types.UnifiedOpenOrder, error) {
	var response []OrderStatusResponse
	err := c.Call(http.MethodPost, "spot/v4/query/open-orders", map[string]interface{}{
		"symbol": symbol,
	}, constants.Signed, &response)
	if err != nil {
		return nil, err
	}
	var result = make([]types.UnifiedOpenOrder, 0, len(response))

	for _, order := range response {
		price, _ := decimal.NewFromString(order.Price)
		amount, _ := decimal.NewFromString(order.Size)
		result = append(result, types.UnifiedOpenOrder{
			OrderId:  order.OrderID,
			TradeNo:  order.ClientOrderID,
			Symbol:   order.Symbol,
			Side:     strings.ToUpper(order.Side),
			Type:     OrderType(order.Type).Convert(),
			Price:    price,
			Quantity: amount,
			Status:   OrderStatus(order.State).Convert(),
		})
	}
	return result, nil
}

type OrderBookResponse struct {
	Timestamp string     `json:"ts"`     // Create time(Timestamp in milliseconds)
	Symbol    string     `json:"symbol"` // Trading pair
	Asks      [][]string `json:"asks"`   // Order book on sell side
	Bids      [][]string `json:"bids"`   // Order book on buy side
	Amount    string     `json:"amount"` // Total number of current price depth
	Price     string     `json:"price"`  // The price at current depth
}

func (c *Connector) OrderBook(symbol string, limit *int64) (*types.UnifiedOrderBook, error) {
	var response OrderBookResponse
	if limit == nil {
		*limit = 30
	}
	err := c.Call(http.MethodGet, "spot/quotation/v3/books", map[string]interface{}{
		"symbol": symbol,
		"limit":  limit,
	}, constants.None, &response)
	if err != nil {
		return nil, err
	}
	ts, err := strconv.ParseInt(response.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	return &types.UnifiedOrderBook{
		Symbol:    symbol,
		Asks:      response.Asks,
		Bids:      response.Bids,
		Timestamp: ts,
	}, nil
}
