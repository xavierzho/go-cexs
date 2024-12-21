package binance

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

// NewOrderACK ACK
type NewOrderACK struct {
	Symbol        string `json:"symbol"`
	OrderId       int64  `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
	IsIsolated    bool   `json:"isIsolated"`
	TransactTime  uint64 `json:"transactTime"`
}
type NewOrderRESULT struct {
	NewOrderACK
	Price                   string `json:"price,omitempty"`
	OrigQty                 string `json:"origQty,omitempty"`
	ExecutedQty             string `json:"executedQty,omitempty"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty,omitempty"`
	Status                  string `json:"status,omitempty"`
	TimeInForce             string `json:"timeInForce,omitempty"`
	Type                    string `json:"type,omitempty"`
	Side                    string `json:"side,omitempty"`
	WorkingTime             uint64 `json:"workingTime,omitempty"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode,omitempty"`
}

type NewOrderFULL struct {
	NewOrderRESULT
	MarginBuyBorrowAmount float64 `json:"marginBuyBorrowAmount,omitempty"`
	MarginBuyBorrowAsset  string  `json:"marginBuyBorrowAsset,omitempty"`
	Fills                 []struct {
		Price           string `json:"price"`
		Qty             string `json:"qty"`
		Commission      string `json:"commission"`
		CommissionAsset string `json:"commissionAsset"`
	} `json:"fills,omitempty"`
}

func (c *Connector) MatchOrderType(orderType constants.OrderType) types.OrderTypeConverter {
	switch orderType {
	case constants.Market:
		return OrderTypeMarket
	case constants.Limit:
		return OrderTypeLimit
	case constants.LimitMaker:
		return OrderTypeLimitMaker
	case constants.StopLossLimit:
		return OrderTypeStopLossLimit
	case constants.StopLoss:
		return OrderTypeStopLoss
	default:
		return OrderTypeMarket
	}
}

func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	resp := new(NewOrderFULL)
	//fmt.Println("request", params)
	orderType := c.MatchOrderType(params.Type)
	err := c.Call(http.MethodPost, OrderEndpoint, map[string]any{
		SymbolFiled:        params.Symbol,
		"side":             strings.ToUpper(params.Side),
		"type":             orderType,
		"price":            params.Price.StringFixed(8),
		"quantity":         params.Quantity.StringFixed(2),
		"timeInForce":      GTC.String(),
		"newClientOrderId": params.TradeNo,
		"newOrderRespType": NewOrderRespTypeFULL,
	}, constants.Signed, resp)

	return strconv.FormatInt(resp.OrderId, 10), err
}

func (c *Connector) BatchOrder(orders []types.OrderEntry) ([]string, error) {
	var list []string
	for _, order := range orders {
		orderId, err := c.PlaceOrder(order)
		if err != nil {
			continue
		}
		list = append(list, orderId)
	}
	return list, nil
}

type QueryOrder struct {
	Symbol                  string `json:"symbol"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Side                    string `json:"side"`
	OrderListID             int    `json:"orderListId"`
	ExecutedQty             string `json:"executedQty"`
	OrderID                 int    `json:"orderId"`
	OrigQty                 string `json:"origQty"`
	ClientOrderID           string `json:"clientOrderId"`
	UpdateTime              int64  `json:"updateTime"`
	WorkingTime             int64  `json:"workingTime"`
	Type                    string `json:"type"`
	IcebergQty              string `json:"icebergQty"`
	StopPrice               string `json:"stopPrice"`
	Price                   string `json:"price"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty"`
	Time                    int64  `json:"time"`
	TimeInForce             string `json:"timeInForce"`
	IsWorking               bool   `json:"isWorking"`
	Status                  string `json:"status"`
}
type CancelOrder struct {
	Symbol              string `json:"symbol"`
	IsIsolated          bool   `json:"isIsolated"`
	OrderId             int64  `json:"orderId"`
	OrigClientOrderId   string `json:"origClientOrderId"`
	ClientOrderId       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
}

func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	var resp = new(CancelOrder)
	od, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return false, err
	}
	err = c.Call(http.MethodDelete, OrderEndpoint, map[string]any{
		SymbolFiled: symbol,
		"orderId":   od,
	}, constants.Signed, resp)
	if err != nil {
		return false, err
	}
	if resp.OrderId != od {
		return false, nil
	}
	//fmt.Println("canceled", resp)
	return true, nil
}
func (c *Connector) CancelAll(symbol string) error {
	return c.Call(http.MethodDelete, OpenOrdersEndpoint, map[string]any{
		SymbolFiled: symbol,
	}, constants.Signed, nil)
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var result = make(map[string]bool)
	for _, id := range orderIds {
		success, err := c.Cancel(symbol, id)
		if err != nil {
			continue
		}
		result[id] = success
	}
	return result, nil
}

func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {
	var resp = new(QueryOrder)
	err := c.Call(http.MethodGet, OrderEndpoint, map[string]any{
		SymbolFiled: symbol,
		"orderId":   orderId,
	}, constants.Signed, resp)
	if err != nil {
		return constants.Error, err
	}
	return OrderStatus(resp.Status).Convert(), err
}

type OpenOrder struct {
	Symbol                  string `json:"symbol"`
	CummulativeQuoteQty     string `json:"cummulativeQuoteQty"`
	Side                    string `json:"side"`
	OrderListID             int64  `json:"orderListId"`
	ExecutedQty             string `json:"executedQty"`
	OrderID                 int64  `json:"orderId"`
	OrigQty                 string `json:"origQty"`
	ClientOrderID           string `json:"clientOrderId"`
	UpdateTime              int64  `json:"updateTime"`
	WorkingTime             int64  `json:"workingTime"`
	Type                    string `json:"type"`
	IcebergQty              string `json:"icebergQty"`
	StopPrice               string `json:"stopPrice"`
	Price                   string `json:"price"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	OrigQuoteOrderQty       string `json:"origQuoteOrderQty"`
	Time                    int64  `json:"time"`
	TimeInForce             string `json:"timeInForce"`
	IsWorking               bool   `json:"isWorking"`
	Status                  string `json:"status"`
}

func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var openOrders []OpenOrder
	err := c.Call(http.MethodGet, OpenOrdersEndpoint, map[string]any{
		SymbolFiled: symbol,
	}, constants.Signed, &openOrders)
	if err != nil {
		return nil, err
	}
	var result []types.OpenOrderEntry
	for _, order := range openOrders {
		price, _ := decimal.NewFromString(order.Price)
		quantity, _ := decimal.NewFromString(order.OrigQty)
		result = append(result, types.OpenOrderEntry{
			OrderId:  strconv.FormatInt(order.OrderID, 10),
			Type:     OrderType(order.Type).Convert(),
			Side:     order.Side,
			Price:    price,
			Quantity: quantity,
			Status:   OrderStatus(order.Status).Convert(),
			Symbol:   order.Symbol,
			TradeNo:  order.ClientOrderID,
		})
	}
	return result, nil
}
