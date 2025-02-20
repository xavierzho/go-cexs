package mexc

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Order struct {
	Symbol              string `json:"symbol"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Side                string `json:"side"`
	OrderListID         int    `json:"orderListId"`
	ExecutedQty         string `json:"executedQty"`
	OrderID             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	ClientOrderID       string `json:"clientOrderId"`
	Type                string `json:"type"`
	OrigClientOrderID   string `json:"origClientOrderId"`
	Price               string `json:"price"`
	TimeInForce         string `json:"timeInForce"`
	Status              string `json:"status"`
}

func (c *Connector) MatchOrderType(orderType constants.OrderType) types.OrderTypeConverter {
	switch orderType {
	case constants.Market:
		return OrderTypeMarket
	case constants.Limit:
		return OrderTypeLimit
	case constants.LimitMaker:
		return OrderTypeLimitMarker
	default:
		return OrderTypeMarket
	}
}
func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	var resp = new(Order)

	orderType := c.MatchOrderType(params.Type)
	err := c.Call(http.MethodPost, OrderEndpoint, &platforms.ObjectBody{
		SymbolFiled: params.Symbol,
		"side":      strings.ToUpper(params.Side),
		"type":      orderType,
	}, constants.Signed, resp)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(int64(resp.OrderID), 10), nil
}

func (c *Connector) BatchOrder(params []types.OrderEntry) ([]string, error) {

	orders := make(platforms.ArrayBody, len(params))
	for i, arg := range params {
		orders[i] = map[string]interface{}{
			"quantity":         arg.Quantity.StringFixed(1),
			"price":            arg.Price.StringFixed(11),
			"side":             strings.ToLower(arg.Side),
			SymbolFiled:        arg.Symbol,
			"type":             c.MatchOrderType(arg.Type).String(),
			"newClientOrderId": arg.TradeNo,
		}
	}
	maxSize := 30
	numBatches := int(math.Ceil(float64(len(orders)) / float64(maxSize)))
	results := make([]string, len(orders))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < numBatches; i++ {
		wg.Add(1)
		batchStart := i * maxSize
		batchEnd := int(math.Min(float64(batchStart+maxSize), float64(len(orders))))
		batchOrders := orders[batchStart:batchEnd]
		go func(i int, batchOrders []map[string]interface{}) {
			defer wg.Done()
			var resp []Order
			err := c.Call(http.MethodPost, BatchOrderEndpoint, &platforms.ObjectBody{
				"batchOrders": batchOrders,
			}, constants.Signed, &resp)
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			if len(resp) > 0 {
				for j, order := range resp {
					results[batchStart+j] = strconv.FormatInt(int64(order.OrderID), 10)
				}
			}
		}(i, batchOrders)
	}

	return results, nil
}
func (c *Connector) queryOrder(symbol string, orderId string) (types.QueryOrder, error) {
	var resp types.QueryOrder
	err := c.Call(http.MethodGet, OrderEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		orderId:     orderId,
	}, constants.Signed, &resp)
	if err != nil {
		return types.QueryOrder{}, err
	}
	return resp, nil
}
func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {
	order, err := c.queryOrder(symbol, orderId)
	if err != nil {
		return constants.Error, err
	}
	return order.Status, nil
}
func (c *Connector) QueryOrder(symbol string, orderId string) (types.QueryOrder, error) {
	return c.queryOrder(symbol, orderId)
}
func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	err := c.Call(http.MethodDelete, OrderEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		orderId:     orderId,
	}, constants.Signed, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Connector) CancelAll(symbol string) error {
	return c.Call(http.MethodDelete, OpenOrdersEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, nil)
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var wg sync.WaitGroup
	var mux sync.Mutex
	var result = make(map[string]bool)
	for _, orderId := range orderIds {
		wg.Add(1)
		go func(orderId string) {
			defer wg.Done()

			ok, err := c.Cancel(symbol, orderId)
			if err != nil {
				return
			}
			mux.Lock()
			result[orderId] = ok
			mux.Unlock()
		}(orderId)
	}
	wg.Wait()
	return result, nil
}

type OpenOrder struct {
	Symbol              string `json:"symbol"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Side                string `json:"side"`
	OrderListID         int    `json:"orderListId"`
	ExecutedQty         string `json:"executedQty"`
	OrderID             int    `json:"orderId"`
	OrigQty             string `json:"origQty"`
	ClientOrderID       string `json:"clientOrderId"`
	UpdateTime          int64  `json:"updateTime"`
	Type                string `json:"type"`
	IcebergQty          string `json:"icebergQty"`
	StopPrice           string `json:"stopPrice"`
	Price               string `json:"price"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
	Time                int64  `json:"time"`
	TimeInForce         string `json:"timeInForce"`
	IsWorking           bool   `json:"isWorking"`
	Status              string `json:"status"`
}

func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var resp []OpenOrder
	var result []types.OpenOrderEntry
	err := c.Call(http.MethodGet, OpenOrdersEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	for _, m := range resp {
		price, _ := decimal.NewFromString(m.Price)
		amount, _ := decimal.NewFromString(m.OrigQty)
		result = append(result, types.OpenOrderEntry{
			Symbol:   m.Symbol,
			Type:     OrderType(m.Type).Convert(),
			Status:   OrderStatus(m.Status).Convert(),
			Side:     m.Side,
			OrderId:  strconv.Itoa(m.OrderID),
			TradeNo:  m.ClientOrderID,
			Price:    price,
			Quantity: amount,
		})
	}
	return result, nil
}

type Balance struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
	CanWithdraw bool        `json:"canWithdraw"`
	Permissions []string    `json:"permissions"`
	AccountType string      `json:"accountType"`
	UpdateTime  interface{} `json:"updateTime"`
	CanDeposit  bool        `json:"canDeposit"`
	CanTrade    bool        `json:"canTrade"`
}

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var resp Balance
	var result map[string]types.BalanceEntry
	err := c.Call(http.MethodGet, AccountEndpoint, &platforms.ObjectBody{}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	for _, balance := range resp.Balances {
		if len(symbols) > 0 && !slices.Contains(symbols, balance.Asset) {
			continue
		}
		result[balance.Asset] = types.BalanceEntry{
			Currency: balance.Asset,
			Locked:   balance.Locked,
			Free:     balance.Free,
		}
	}

	return result, nil
}
