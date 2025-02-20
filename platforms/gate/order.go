package gate

import (
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

type Order struct {
	Fee                string `json:"fee,omitempty"`
	FeeCurrency        string `json:"fee_currency,omitempty"`
	RebatedFee         string `json:"rebated_fee,omitempty"`
	Type               string `json:"type,omitempty"`
	AutoRepay          bool   `json:"auto_repay,omitempty"`
	FilledTotal        string `json:"filled_total,omitempty"`
	RebatedFeeCurrency string `json:"rebated_fee_currency,omitempty"`
	UpdateTime         string `json:"update_time,omitempty"`
	PointFee           string `json:"point_fee,omitempty"`
	CurrencyPair       string `json:"currency_pair"`
	Price              string `json:"price,omitempty"`
	Iceberg            string `json:"iceberg,omitempty"`
	ID                 string `json:"id,omitempty"`
	Text               string `json:"text,omitempty"`
	AutoBorrow         bool   `json:"auto_borrow,omitempty"`
	UpdateTimeMs       int    `json:"update_time_ms,omitempty"`
	GtMakerFee         string `json:"gt_maker_fee,omitempty"`
	FinishAs           string `json:"finish_as,omitempty"`
	Side               string `json:"side"`
	Amount             string `json:"amount"`
	CreateTime         string `json:"create_time,omitempty"`
	AvgDealPrice       string `json:"avg_deal_price,omitempty"`
	GtFee              string `json:"gt_fee,omitempty"`
	GtDiscount         bool   `json:"gt_discount,omitempty"`
	AmendText          string `json:"amend_text,omitempty"`
	ActionMode         string `json:"action_mode,omitempty"`
	TimeInForce        string `json:"time_in_force,omitempty"`
	CreateTimeMs       int    `json:"create_time_ms,omitempty"`
	Left               string `json:"left,omitempty"`
	GtTakerFee         string `json:"gt_taker_fee,omitempty"`
	StpAct             string `json:"stp_act,omitempty"`
	FilledAmount       string `json:"filled_amount,omitempty"`
	FillPrice          string `json:"fill_price,omitempty"`
	Account            string `json:"account,omitempty"`
	Status             string `json:"status,omitempty"`
	StpID              int    `json:"stp_id,omitempty"`
	OrderId            string `json:"order_id,omitempty"`
	Role               string `json:"role,omitempty"`
}

func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	var resp Order
	var param = &platforms.ObjectBody{
		SymbolFiled:     params.Symbol,
		"text":          fmt.Sprintf("t-%s", uuid.New().String()),
		"type":          c.MatchOrderType(params.Type),
		"side":          strings.ToLower(params.Side),
		"amount":        params.Quantity.StringFixed(10),
		"price":         params.Price.StringFixed(10),
		"time_in_force": TimeInForceGTC,
	}
	if params.Type == constants.Iceberg {
		param.Set("iceberg", params.Quantity.StringFixed(10))
		//param["iceberg"] = params.Quantity.StringFixed(10)
	}
	err := c.Call(http.MethodPost, OrderEndpoint, param, constants.Signed, &resp)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
func (c *Connector) MatchOrderType(orderType constants.OrderType) OrderType {
	switch orderType {
	case constants.Market:
		return OrderTypeMarket
	case constants.Limit:
		return OrderTypeLimit
	default:
		return OrderTypeMarket
	}
}
func (c *Connector) BatchOrder(orders []types.OrderEntry) ([]string, error) {
	var params platforms.ArrayBody
	var result []string
	var resp []Order
	for _, order := range orders {
		params = append(params, map[string]any{
			SymbolFiled:     order.Symbol,
			"text":          fmt.Sprintf("t-%s", uuid.New().String()),
			"type":          c.MatchOrderType(order.Type),
			"side":          strings.ToLower(order.Side),
			"amount":        order.Quantity.StringFixed(10),
			"price":         order.Price.StringFixed(10),
			"time_in_force": TimeInForceGTC,
		})
	}
	err := c.Call(http.MethodPost, BatchOrdersEndpoint, &params, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	for _, r := range resp {
		result = append(result, r.OrderId)
	}
	return result, nil
}

func (c *Connector) queryOrder(symbol string, orderId string) (Order, error) {
	var resp Order
	err := c.Call(http.MethodGet, fmt.Sprintf("%s/%s", OrderEndpoint, orderId), &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, &resp)
	if err != nil {
		return Order{}, err
	}
	return resp, nil
}
func (c *Connector) QueryOrder(symbol string, orderId string) (types.QueryOrder, error) {
	order, err := c.queryOrder(symbol, orderId)
	if err != nil {
		return types.QueryOrder{}, err
	}
	price, _ := decimal.NewFromString(order.Price)
	amount, _ := decimal.NewFromString(order.Amount)
	filled, _ := decimal.NewFromString(order.FilledAmount)
	return types.QueryOrder{
		Symbol:     symbol,
		Type:       OrderType(order.Type).Convert(),
		Status:     OrderStatus(order.Status).Convert(),
		Side:       strings.ToUpper(order.Side),
		Price:      price,
		Quantity:   amount,
		Filled:     filled,
		CreateTime: int64(order.CreateTimeMs / 1000),
		UpdateTime: int64(order.UpdateTimeMs / 1000),
		OrderId:    orderId,
		TradeNo:    order.Text,
	}, nil
}
func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {
	order, err := c.queryOrder(symbol, orderId)
	if err != nil {
		return constants.Error, err
	}
	return OrderStatus(order.Status).Convert(), nil
}

func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	var resp Order
	err := c.Call(http.MethodDelete, fmt.Sprintf("%s/%s", OrderEndpoint, orderId), &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, &resp)
	if err != nil {
		return false, err
	}
	return resp.Status == OrderStatusClosed.String(), nil
}

func (c *Connector) CancelAll(symbol string) error {
	var resp []Order
	return c.Call(http.MethodDelete, OrderEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, &resp)
}

type CancelById struct {
	CurrencyPair string `json:"currency_pair"`
	ID           string `json:"id"`
	Text         string `json:"text"`
	Label        string `json:"label"`
	Message      string `json:"message"`
	Succeeded    bool   `json:"succeeded"`
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var body []map[string]any

	for _, id := range orderIds {
		body = append(body, map[string]any{
			SymbolFiled: symbol,
			"order_id":  id,
		})
	}
	var resp []CancelById
	err := c.Call(http.MethodPost, BatchCancelEndpoint, &platforms.ObjectBody{
		"orders": body,
	}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]bool)
	for _, cancel := range resp {
		result[cancel.ID] = cancel.Succeeded
	}
	return result, nil
}

type PendingOrder struct {
	Symbol string  `json:"currency_pair"`
	Total  int     `json:"total"`
	Orders []Order `json:"orders"`
}

func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var resp []PendingOrder
	err := c.Call(http.MethodGet, OpenOrdersEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
	}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	var result []types.OpenOrderEntry
	for _, order := range resp {
		for _, o := range order.Orders {
			price, _ := decimal.NewFromString(o.Price)
			quantity, _ := decimal.NewFromString(o.Amount)
			result = append(result, types.OpenOrderEntry{
				Symbol:   order.Symbol,
				Type:     OrderType(o.Type).Convert(),
				Status:   OrderStatus(o.Status).Convert(),
				Side:     strings.ToUpper(o.Side),
				OrderId:  o.ID,
				Price:    price,
				Quantity: quantity,
			})
		}
	}
	return result, nil
}

type Balance struct {
	Currency         string `json:"currency"`
	AvailableBalance string `json:"available_balance"`
	EstimatedAsBtc   string `json:"estimated_as_btc"`
	ConvertibleToGt  string `json:"convertible_to_gt"`
}

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var resp [][]Balance
	err := c.Call(http.MethodGet, SmallBalanceEndpoint, &platforms.ObjectBody{}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	var result = make(map[string]types.BalanceEntry)
	for _, balances := range resp {
		for _, balance := range balances {
			result[balance.Currency] = types.BalanceEntry{
				Free:     balance.AvailableBalance,
				Currency: balance.Currency,
			}
		}
	}
	return result, nil
}
