package bybit

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"net/http"
	"strconv"
	"strings"
)

type Order struct {
	OrderId string `json:"orderId"`
	TradeNo string `json:"orderLinkId"`
}

func (o Order) String() string {
	return o.OrderId
}

func FirstSide(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func (c *Connector) MatchOrderType(orderType constants.OrderType) types.OrderTypeConverter {
	switch orderType {
	case constants.Limit:
		return OrderTypeLimit
	case constants.Market:
		return OrderTypeMarket
	default:
		return OrderTypeMarket
	}
}

func (c *Connector) RawPlaceOrder(params platforms.Serializer) (Order, error) {
	var resp RestResp[Order, NullExt]
	err := c.Call(http.MethodPost, PlaceOrderEndpoint, params, constants.Signed, &resp)
	if err != nil {
		return Order{}, err
	}
	return resp.Result, nil
}

func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	var p = &platforms.ObjectBody{
		"category":   "spot",
		"symbol":     params.Symbol,
		"isLeverage": false,
		"side":       FirstSide(params.Side),
		"orderType":  c.MatchOrderType(params.Type).String(),
		"qty":        params.Quantity.StringFixed(2),
	}

	if !params.Price.IsZero() {
		p.Set("price", params.Price.StringFixed(12))
	}
	order, err := c.RawPlaceOrder(p)
	if err != nil {
		return "", err
	}
	return order.OrderId, nil
}

type OrderList struct {
	List []struct {
		Symbol      string `json:"symbol"`
		OrderLinkID string `json:"orderLinkId"`
		OrderID     string `json:"orderId"`
		Category    string `json:"category"`
		CreateAt    string `json:"createAt,omitempty"`
	} `json:"list"`
}

func (OrderList) String() string {
	return ""
}

type OrdersExt struct {
	List []struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	} `json:"list"`
}

func (ext OrdersExt) String() string {
	return ""
}
func (c *Connector) RawBatchOrder(orders []map[string]any) (OrderList, error) {
	var resp RestResp[OrderList, OrdersExt]

	err := c.Call(http.MethodPost, BatchPlaceOrderEndpoint, &platforms.ObjectBody{
		"category": orders[0]["category"],
		"request":  orders,
	}, constants.Signed, &resp)
	if err != nil {
		return OrderList{}, err
	}
	return resp.Result, nil
}
func (c *Connector) BatchOrder(orders []types.OrderEntry) ([]string, error) {
	var result []string
	var params = make([]map[string]any, len(orders))
	for i, order := range orders {
		params[i] = map[string]any{
			"category":    "spot",
			"symbol":      order.Symbol,
			"orderType":   c.MatchOrderType(order.Type).String(),
			"side":        FirstSide(order.Side),
			"qty":         order.Quantity.StringFixed(2),
			"price":       order.Price.StringFixed(12),
			"timeInForce": "GTC",
			"orderLinkId": uuid.New().String(),
			"mmp":         false,
			"reduceOnly":  false,
		}
	}
	resp, err := c.RawBatchOrder(params)
	if err != nil {
		return nil, err
	}
	for _, ord := range resp.List {
		result = append(result, ord.OrderID)
	}
	return result, nil
}

type OrderInfo struct {
	OrderId            string `json:"orderId"`
	OrderLinkId        string `json:"orderLinkId"`
	BlockTradeId       string `json:"blockTradeId"`
	Symbol             string `json:"symbol"`
	Price              string `json:"price"`
	Qty                string `json:"qty"`
	Side               string `json:"side"`
	IsLeverage         string `json:"isLeverage"`
	PositionIdx        int    `json:"positionIdx"`
	OrderStatus        string `json:"orderStatus"`
	CancelType         string `json:"cancelType"`
	RejectReason       string `json:"rejectReason"`
	AvgPrice           string `json:"avgPrice"`
	LeavesQty          string `json:"leavesQty"`
	LeavesValue        string `json:"leavesValue"`
	CumExecQty         string `json:"cumExecQty"`
	CumExecValue       string `json:"cumExecValue"`
	CumExecFee         string `json:"cumExecFee"`
	TimeInForce        string `json:"timeInForce"`
	OrderType          string `json:"orderType"`
	StopOrderType      string `json:"stopOrderType"`
	OrderIv            string `json:"orderIv"`
	TriggerPrice       string `json:"triggerPrice"`
	TakeProfit         string `json:"takeProfit"`
	StopLoss           string `json:"stopLoss"`
	TpslMode           string `json:"tpslMode"`
	OcoTriggerType     string `json:"ocoTriggerType"`
	TpLimitPrice       string `json:"tpLimitPrice"`
	SlLimitPrice       string `json:"slLimitPrice"`
	TpTriggerBy        string `json:"tpTriggerBy"`
	SlTriggerBy        string `json:"slTriggerBy"`
	TriggerDirection   int    `json:"triggerDirection"`
	TriggerBy          string `json:"triggerBy"`
	LastPriceOnCreated string `json:"lastPriceOnCreated"`
	ReduceOnly         bool   `json:"reduceOnly"`
	CloseOnTrigger     bool   `json:"closeOnTrigger"`
	PlaceType          string `json:"placeType"`
	SmpType            string `json:"smpType"`
	SmpGroup           int    `json:"smpGroup"`
	SmpOrderId         string `json:"smpOrderId"`
	CreatedTime        string `json:"createdTime"`
	UpdatedTime        string `json:"updatedTime"`
}

type OrderInfos struct {
	List           []OrderInfo `json:"list"`
	NextPageCursor string      `json:"nextPageCursor"`
	Category       string      `json:"category"`
}

func (OrderInfos) String() string {
	return ""
}
func (c *Connector) RawOrder(params platforms.Serializer) (OrderInfo, error) {
	var resp RestResp[OrderInfos, NullExt]
	err := c.Call(http.MethodGet, RealTimeOrderEndpoint, params, constants.Signed, &resp)
	if err != nil {
		return OrderInfo{}, err
	}
	return resp.Result.List[0], nil
}
func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {

	order, err := c.RawOrder(&platforms.ObjectBody{
		"symbol":  symbol,
		"orderId": orderId,
	})
	if err != nil {
		return constants.Error, err
	}

	return OrderStatus(order.OrderStatus).Convert(), nil
}
func (c *Connector) QueryOrder(symbol string, orderId string) (types.QueryOrder, error) {
	var result types.QueryOrder
	order, err := c.RawOrder(&platforms.ObjectBody{
		"symbol":  symbol,
		"orderId": orderId,
	})
	if err != nil {
		return result, err
	}

	result.OrderId = order.OrderId
	result.Symbol = order.Symbol
	result.Side = strings.ToUpper(order.Side)
	qty, _ := decimal.NewFromString(order.Qty)
	result.Quantity = qty
	price, _ := decimal.NewFromString(order.Price)
	result.Price = price
	executed, _ := decimal.NewFromString(order.CumExecQty)
	result.Filled = qty.Sub(executed)
	result.TradeNo = order.OrderLinkId
	result.Status = OrderStatus(order.OrderStatus).Convert()
	created, _ := strconv.ParseInt(order.CreatedTime, 10, 64)
	result.CreateTime = created
	updated, _ := strconv.ParseInt(order.UpdatedTime, 10, 16)
	result.UpdateTime = updated
	return result, nil
}
func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	var resp RestResp[Order, NullExt]
	err := c.Call(http.MethodPost, OrderCancelEndpoint, &platforms.ObjectBody{
		"symbol":   symbol,
		"orderId":  orderId,
		"category": "spot",
	}, constants.Signed, &resp)
	if err != nil {
		return false, err
	}
	return true, nil
}

type Orders []Order

func (Orders) String() string {
	return ""
}

func (c *Connector) CancelAll(symbol string) error {
	var resp RestResp[Orders, NullExt]
	return c.Call(http.MethodPost, OrderCancelAllEndpoint, &platforms.ObjectBody{
		"symbol":   symbol,
		"category": "spot",
	}, constants.Signed, &resp)
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var resp RestResp[OrderList, OrdersExt]
	var req = make([]map[string]string, len(orderIds))
	var result = make(map[string]bool)
	for i, id := range orderIds {
		req[i] = map[string]string{
			"symbol":  symbol,
			"orderId": id,
		}
	}
	err := c.Call(http.MethodPost, OrderBatchCancelEndpoint, &platforms.ObjectBody{
		"category": "spot",
		"request":  req,
	}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}
	for i, info := range resp.Ext.List {
		result[resp.Result.List[i].OrderID] = info.Code == 0 || info.Msg == "OK"
	}
	return result, nil
}
func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var result []types.OpenOrderEntry
	var params = &platforms.ObjectBody{
		"category": "spot",
		"symbol":   symbol,
		"limit":    50,
	}

	for {
		var resp RestResp[OrderInfos, NullExt]
		err := c.Call(http.MethodGet, RealTimeOrderEndpoint, params, constants.Signed, &resp)
		if err != nil {
			return nil, err
		}
		for _, info := range resp.Result.List {
			result = append(result, types.OpenOrderEntry{
				Symbol: info.Symbol,
				Side:   strings.ToUpper(info.Side),
			})
		}
		if len(resp.Result.List) != 50 {
			break
		}
		params.Set("cursor", resp.Result.NextPageCursor)
	}

	return result, nil
}

type WalletAccountInfo struct {
	AccountType            string     `json:"accountType"`
	AccountLTV             string     `json:"accountLTV"`
	AccountIMRate          string     `json:"accountIMRate"`
	AccountMMRate          string     `json:"accountMMRate"`
	TotalEquity            string     `json:"totalEquity"`
	TotalWalletBalance     string     `json:"totalWalletBalance"`
	TotalMarginBalance     string     `json:"totalMarginBalance"`
	TotalAvailableBalance  string     `json:"totalAvailableBalance"`
	TotalPerpUPL           string     `json:"totalPerpUPL"`
	TotalInitialMargin     string     `json:"totalInitialMargin"`
	TotalMaintenanceMargin string     `json:"totalMaintenanceMargin"`
	Coins                  []CoinInfo `json:"coin"`
}

type CoinInfo struct {
	Coin                string `json:"coin"`
	Equity              string `json:"equity"`
	UsdValue            string `json:"usdValue"`
	WalletBalance       string `json:"walletBalance"`
	Free                string `json:"free"`
	Locked              string `json:"locked"`
	BorrowAmount        string `json:"borrowAmount"`
	AvailableToBorrow   string `json:"availableToBorrow"`
	AvailableToWithdraw string `json:"availableToWithdraw"`
	AccruedInterest     string `json:"accruedInterest"`
	TotalOrderIM        string `json:"totalOrderIM"`
	TotalPositionIM     string `json:"totalPositionIM"`
	TotalPositionMM     string `json:"totalPositionMM"`
	UnrealisedPnl       string `json:"unrealisedPnl"`
	CumRealisedPnl      string `json:"cumRealisedPnl"`
	Bonus               string `json:"bonus"`
	CollateralSwitch    bool   `json:"collateralSwitch"`
	MarginCollateral    bool   `json:"marginCollateral"`
}

type WalletBalance struct {
	List []WalletAccountInfo `json:"list"`
}

func (WalletBalance) String() string {
	return ""
}

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var resp RestResp[WalletBalance, NullExt]
	var result = make(map[string]types.BalanceEntry)
	err := c.Call(http.MethodGet, WalletBalanceEndpoint, &platforms.ObjectBody{
		"accountType": SpotAccount,
		"coin":        strings.Join(symbols, ","),
	}, constants.Signed, &resp)
	if err != nil {
		return nil, err
	}

	for _, info := range resp.Result.List {
		for _, coin := range info.Coins {
			result[coin.Coin] = types.BalanceEntry{
				Free:     coin.Free,
				Locked:   coin.Locked,
				Currency: coin.Coin,
			}
		}
	}

	return result, nil
}
