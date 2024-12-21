package bitmart

import (
	"github.com/shopspring/decimal"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
)

type OrderResponse struct {
	OrderId string `json:"order_id"`
}

func (c *Connector) MatchOrderType(state constants.OrderType) types.OrderTypeConverter {
	switch state {
	case constants.Market:
		return OrderTypeMarket
	case constants.Limit:
		return OrderTypeLimit
	case constants.LimitMaker:
		return OrderTypeLimitMaker
	default:
		return OrderTypeMarket
	}
}
func (c *Connector) PlaceOrder(order types.OrderEntry) (string, error) {
	params := map[string]any{
		"symbol":        order.Symbol,
		"side":          strings.ToLower(order.Side),
		"type":          c.MatchOrderType(order.Type).String(),
		"timeInForce":   order.TimeInForce,
		"price":         order.Price.StringFixed(11),
		"quantity":      order.Quantity.StringFixed(1),
		"clientOrderId": order.TradeNo,
		"notional":      "",
	}
	var response OrderResponse
	err := c.Call(http.MethodPost, NewOrderEndpoint, params, constants.Signed, &response)
	if err != nil {
		return "", err
	}

	return response.OrderId, nil
}

type BatchResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		OrderIds []string `json:"order_ids"`
	} `json:"data"`
}

func (c *Connector) BatchOrder(params []types.OrderEntry) ([]string, error) {
	// Prepare orders
	orders := make([]map[string]interface{}, len(params))
	for i, arg := range params {
		orders[i] = map[string]interface{}{
			"size":          arg.Quantity.StringFixed(1),
			"price":         arg.Price.StringFixed(11),
			"side":          strings.ToLower(arg.Side),
			"symbol":        arg.Symbol,
			"type":          c.MatchOrderType(arg.Type).String(),
			"clientOrderId": arg.TradeNo,
		}
	}

	maxSize := 10
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
			params := map[string]interface{}{
				"symbol":      batchOrders[0]["symbol"],
				"orderParams": batchOrders,
			}
			var response BatchResponse
			err := c.Call(http.MethodPost, BatchOrderEndpoint, params, constants.Signed, &response)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				log.Printf("Error in batch %d: %v\n", i, err)
			} else {
				// 将批次的订单 ID 填充到结果中，保证原始顺序
				for j, id := range response.Data.OrderIds {
					results[batchStart+j] = id
				}
			}
		}(i, batchOrders)
	}

	wg.Wait()
	return results, nil
}

type OrderStatusResponse struct {
	Symbol         string `json:"symbol"`
	Side           string `json:"side"`
	Notional       string `json:"notional"`
	OrderID        string `json:"orderId"`
	ClientOrderID  string `json:"clientOrderId"`
	UpdateTime     int64  `json:"updateTime"`
	Type           string `json:"type"`
	PriceAvg       string `json:"priceAvg"`
	OrderMode      string `json:"orderMode"`
	Size           string `json:"size"`
	FilledSize     string `json:"filledSize"`
	CreateTime     int64  `json:"createTime"`
	Price          string `json:"price"`
	FilledNotional string `json:"filledNotional"`
	State          string `json:"state"`
}

func (c *Connector) GetOrderStatus(_ string, orderId string) (constants.OrderStatus, error) {
	var response OrderStatusResponse
	err := c.Call(http.MethodPost, QueryOrderEndpoint, map[string]interface{}{
		"order_id": orderId,
	}, constants.Signed, &response)
	if err != nil {
		return constants.Error, err
	}

	var state = OrderStatus(response.State).Convert()
	return state, nil
}

func (c *Connector) PendingOrders(symbol string) ([]types.OpenOrderEntry, error) {
	var response []OrderStatusResponse
	err := c.Call(http.MethodPost, OpenOrdersEndpoint, map[string]interface{}{
		"symbol": symbol,
	}, constants.Signed, &response)
	if err != nil {
		return nil, err
	}
	var result = make([]types.OpenOrderEntry, 0, len(response))

	for _, order := range response {
		price, _ := decimal.NewFromString(order.Price)
		amount, _ := decimal.NewFromString(order.Size)
		result = append(result, types.OpenOrderEntry{
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
