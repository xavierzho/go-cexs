package okx

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type OrderReturn struct {
	TradeNo string `json:"clOrdId"`
	SCode   string `json:"sCode"`
	Tag     string `json:"tag,omitempty"`
	SMsg    string `json:"sMsg"`
	OrderId string `json:"ordId"`
	Ts      string `json:"ts"`
}

func (r OrderReturn) String() string {
	return r.SMsg
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
func (c *Connector) PlaceOrder(params types.OrderEntry) (string, error) {
	var resp RestReturn[OrderReturn]
	err := c.Call(http.MethodPost, OrderEndpoint, &platforms.ObjectBody{
		"instId":  c.SymbolPattern(params.Symbol),
		"tdMode":  CashMode,
		"clOrdId": uuid.New().String(),
		"side":    strings.ToLower(params.Side),
		"ordType": c.MatchOrderType(params.Type),
		"px":      params.Price.StringFixed(12),
		"sz":      params.Quantity.StringFixed(2),
	}, constants.None, &resp)
	if err != nil {
		return "", err
	}
	return resp.Data[0].OrderId, nil
}

func (c *Connector) BatchOrder(params []types.OrderEntry) ([]string, error) {
	var orders = make(platforms.ArrayBody, len(params))
	for i, order := range params {
		orders[i] = map[string]any{
			"instId":  c.SymbolPattern(order.Symbol),
			"tdMode":  CashMode,
			"clOrdId": uuid.New().String(),
			"side":    strings.ToLower(order.Side),
			"ordType": c.MatchOrderType(order.Type),
			"px":      order.Price.StringFixed(12),
			"sz":      order.Quantity.StringFixed(2),
		}
	}
	results := make([]string, len(orders))
	const maxOrders = 20
	numBatches := int(math.Ceil(float64(len(orders)) / float64(maxOrders)))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numBatches; i++ {
		wg.Add(1)
		batchStart := i * maxOrders
		batchEnd := int(math.Min(float64(batchStart+maxOrders), float64(len(orders))))
		batchOrders := orders[batchStart:batchEnd]
		go func(i int, batchOrders []map[string]interface{}) {
			defer wg.Done()
			var resp RestReturn[OrderReturn]
			err := c.Call(http.MethodPost, OrderBatchEndpoint, &orders, constants.None, &resp)
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				log.Printf("Error in batch %d: %v\n", i, err)
			} else {
				// 将批次的订单 ID 填充到结果中，保证原始顺序
				for j, order := range resp.Data {
					results[batchStart+j] = order.OrderId
				}
			}
		}(i, batchOrders)
	}
	wg.Wait()

	return results, nil
}

type OrderInfo struct {
	PxType            string `json:"pxType,omitempty"`
	Fee               string `json:"fee"`
	Price             string `json:"px"`
	TpTriggerPxType   string `json:"tpTriggerPxType"`
	Source            string `json:"source"`
	OrderId           string `json:"ordId"`
	AttachAlgoClOrdID string `json:"attachAlgoClOrdId,omitempty"`
	ClientOrderId     string `json:"clOrdId"`
	Ccy               string `json:"ccy,omitempty"` // Margin currency
	LinkedAlgoOrd     struct {
		AlgoID string `json:"algoId"`
	} `json:"linkedAlgoOrd"`
	State              string        `json:"state"`
	Tag                string        `json:"tag"`
	QuickMgnType       string        `json:"quickMgnType"`
	AttachAlgoOrds     []interface{} `json:"attachAlgoOrds,omitempty"`
	SlTriggerPxType    string        `json:"slTriggerPxType"`
	StpID              string        `json:"stpId"`
	Lever              string        `json:"lever"`  // Leverage, from 0.01 to 125. Only applicable to MARGIN/FUTURES/SWAP
	TradeMode          string        `json:"tdMode"` // trade mode
	TgtCcy             string        `json:"tgtCcy"`
	TpOrdPx            string        `json:"tpOrdPx"`
	CancelSourceReason string        `json:"cancelSourceReason"`
	Pnl                string        `json:"pnl"`
	InstType           string        `json:"instType"`
	ReduceOnly         string        `json:"reduceOnly"`
	SlOrdPx            string        `json:"slOrdPx"`
	PriceUSD           string        `json:"pxUsd,omitempty"`
	OrderType          string        `json:"ordType"`
	LatestFillQty      string        `json:"fillSz"`
	PxVol              string        `json:"pxVol,omitempty"`
	AlgoClOrdID        string        `json:"algoClOrdId,omitempty"`
	CreateTime         string        `json:"cTime"`
	TpTriggerPx        string        `json:"tpTriggerPx"`
	AccFilledQty       string        `json:"accFillSz"` // Accumulated fill quantity
	PosSide            string        `json:"posSide"`   // position side
	IsTpLimit          string        `json:"isTpLimit"`
	StpMode            string        `json:"stpMode"`
	Side               string        `json:"side"`
	LatestFillPrice    string        `json:"fillPx"` // Last filled price. If none is filled, it will return "".
	AlgoID             string        `json:"algoId,omitempty"`
	Rebate             string        `json:"rebate"`
	Qty                string        `json:"sz"`
	Symbol             string        `json:"instId"` // Instrument ID, symbol
	AvgPrice           string        `json:"avgPx"`  // Average filled price. If none is filled, it will return "".
	CancelSource       string        `json:"cancelSource"`
	SlTriggerPx        string        `json:"slTriggerPx"`
	UpdateTime         string        `json:"uTime"`
	LatestFillTime     string        `json:"fillTime"`
	Category           string        `json:"category"`
	RebateCcy          string        `json:"rebateCcy"`
	TradeID            string        `json:"tradeId"` // Last traded ID
	FeeCcy             string        `json:"feeCcy"`
}

func (OrderInfo) String() string {
	return ""
}
func (c *Connector) RawOrder(symbol string, orderId string) (OrderInfo, error) {
	var resp RestReturn[OrderInfo]
	symbol = c.SymbolPattern(symbol)
	err := c.Call(http.MethodGet, OrderEndpoint, &platforms.ObjectBody{
		"instId": symbol,
		"ordId":  orderId,
	}, constants.None, &resp)
	if err != nil {
		return OrderInfo{}, err
	}
	return resp.Data[0], nil
}
func (c *Connector) GetOrderStatus(symbol string, orderId string) (constants.OrderStatus, error) {
	order, err := c.RawOrder(symbol, orderId)
	if err != nil {
		return constants.Error, err
	}
	return OrderStatus(order.State).Convert(), nil
}
func (c *Connector) QueryOrder(symbol string, orderId string) (types.QueryOrder, error) {
	order, err := c.RawOrder(symbol, orderId)
	if err != nil {
		return types.QueryOrder{}, err
	}
	ct, _ := strconv.ParseInt(order.CreateTime, 10, 64)
	ut, _ := strconv.ParseInt(order.UpdateTime, 10, 64)
	price, _ := decimal.NewFromString(order.Price)
	qty, _ := decimal.NewFromString(order.Qty)
	return types.QueryOrder{
		Symbol:     order.Symbol,
		Type:       OrderType(order.OrderType).Convert(),
		Status:     OrderStatus(order.State).Convert(),
		Side:       strings.ToUpper(order.Side),
		TradeNo:    order.ClientOrderId,
		OrderId:    order.OrderId,
		CreateTime: ct,
		UpdateTime: ut,
		Price:      price,
		Quantity:   qty,
	}, nil
}
func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	var resp RestReturn[OrderReturn]

	err := c.Call(http.MethodPost, OrderCancelEndpoint, &platforms.ObjectBody{
		"instId": c.SymbolPattern(symbol),
		"ordId":  orderId,
	}, constants.None, &resp)
	if err != nil {
		return false, err
	}
	if resp.Data[0].SCode != "0" {
		return false, fmt.Errorf("not support")
	}
	return true, nil
}

type CancelAll struct {
	Timestamp   string `json:"ts"`
	Tag         string `json:"tag"`
	TriggerTime string `json:"triggerTime"`
}

func (CancelAll) String() string {
	return ""
}
func (c *Connector) CancelAll(_ string) error {
	var resp RestReturn[CancelAll]

	return c.Call(http.MethodPost, OrderCancelAllAfterEndpoint, &platforms.ObjectBody{"timeOut": 0}, constants.None, &resp)
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var resp RestReturn[OrderReturn]
	var orders = make(platforms.ArrayBody, len(orderIds))
	var result = make(map[string]bool)
	for i, id := range orderIds {
		orders[i] = map[string]any{
			"instId": symbol,
			"ordId":  id,
		}
	}
	err := c.Call(http.MethodPost, OrderCancelBatchEndpoint, &orders, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	for _, d := range resp.Data {
		result[d.OrderId] = d.SCode == "0"
	}
	return result, nil
}
func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var results []types.OpenOrderEntry
	var req = &platforms.ObjectBody{"instType": "SPOT", "instId": symbol}
	for {
		var resp RestReturn[OrderInfo]
		err := c.Call(http.MethodPost, OrderPendingEndpoint, req, constants.None, &resp)
		if err != nil {
			return nil, err
		}
		lens := len(resp.Data)
		if lens < 100 {
			break
		}
		for _, order := range resp.Data {
			price, _ := decimal.NewFromString(order.Price)
			qty, _ := decimal.NewFromString(order.Qty)
			results = append(results, types.OpenOrderEntry{
				Symbol:   order.Symbol,
				Type:     OrderType(order.OrderType).Convert(),
				Side:     strings.ToUpper(order.Side),
				TradeNo:  order.ClientOrderId,
				OrderId:  order.OrderId,
				Price:    price,
				Quantity: qty,
				Status:   OrderStatus(order.State).Convert(),
			})
		}
		(*req)["after"] = resp.Data[lens-1].OrderId
	}

	return results, nil
}

type Balance struct {
	Upl        string `json:"upl"`
	BorrowFroz string `json:"borrowFroz"`
	Mmr        string `json:"mmr"`
	AdjEq      string `json:"adjEq"`
	OrdFroz    string `json:"ordFroz"`
	Details    []struct {
		SpotInUseAmt      string `json:"spotInUseAmt"`
		FrozenBal         string `json:"frozenBal"`
		UplLiab           string `json:"uplLiab"`
		Twap              string `json:"twap"`
		CashBal           string `json:"cashBal"`
		SpotUplRatio      string `json:"spotUplRatio"`
		NotionalLever     string `json:"notionalLever"`
		SmtSyncEq         string `json:"smtSyncEq"`
		RewardBal         string `json:"rewardBal"`
		ClSpotInUseAmt    string `json:"clSpotInUseAmt"`
		Imr               string `json:"imr"`
		StgyEq            string `json:"stgyEq"`
		FixedBal          string `json:"fixedBal"`
		SpotUpl           string `json:"spotUpl"`
		Mmr               string `json:"mmr"`
		Interest          string `json:"interest"`
		MaxSpotInUse      string `json:"maxSpotInUse"`
		OpenAvgPx         string `json:"openAvgPx"`
		SpotIsoBal        string `json:"spotIsoBal"`
		Ccy               string `json:"ccy"`
		IsoUpl            string `json:"isoUpl"`
		AvailEq           string `json:"availEq"`
		TotalPnl          string `json:"totalPnl"`
		TotalPnlRatio     string `json:"totalPnlRatio"`
		BorrowFroz        string `json:"borrowFroz"`
		DisEq             string `json:"disEq"`
		IsoLiab           string `json:"isoLiab"`
		Eq                string `json:"eq"`
		IsoEq             string `json:"isoEq"`
		Liab              string `json:"liab"`
		AccAvgPx          string `json:"accAvgPx"`
		OrdFrozen         string `json:"ordFrozen"`
		AvailBal          string `json:"availBal"`
		Upl               string `json:"upl"`
		SpotBal           string `json:"spotBal"`
		CrossLiab         string `json:"crossLiab"`
		SpotCopyTradingEq string `json:"spotCopyTradingEq"`
		MaxLoan           string `json:"maxLoan"`
		EqUsd             string `json:"eqUsd"`
		UTime             string `json:"uTime"`
		MgnRatio          string `json:"mgnRatio"`
	} `json:"details"`
	NotionalUsd string `json:"notionalUsd"`
	UTime       string `json:"uTime"`
	IsoEq       string `json:"isoEq"`
	TotalEq     string `json:"totalEq"`
	Imr         string `json:"imr"`
	MgnRatio    string `json:"mgnRatio"`
}

func (Balance) String() string {
	return ""
}
func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var resp RestReturn[Balance]
	var result = make(map[string]types.BalanceEntry)
	if len(symbols) > 20 {
		symbols = symbols[:20]
	}
	err := c.Call(http.MethodGet, AccountBalanceEndpoint, &platforms.ObjectBody{
		"ccy": strings.Join(symbols, ","),
	}, constants.None, &resp)
	if err != nil {
		return nil, err
	}
	for _, bal := range resp.Data[0].Details {
		result[bal.Ccy] = types.BalanceEntry{
			Free:     bal.AvailBal,
			Locked:   bal.FrozenBal,
			Currency: bal.Ccy,
		}
	}
	return result, nil
}
