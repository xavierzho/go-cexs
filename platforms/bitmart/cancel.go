package bitmart

import (
	"github.com/xavierzho/go-cexs/constants"
	"net/http"
)

func (c *Connector) CancelAll(symbol string) error {
	var response map[string]interface{}
	err := c.Call(http.MethodPost, "spot/v4/cancel_all", map[string]interface{}{
		"symbol": symbol,
	}, constants.Signed, response)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connector) Cancel(symbol, orderId string) (bool, error) {
	var response struct {
		Result bool `json:"result"`
	}
	err := c.Call(http.MethodPost, "spot/v3/cancel_order", map[string]interface{}{
		"symbol":   symbol,
		"order_id": orderId,
	}, constants.Signed, &response)
	if err != nil {
		return false, err
	}
	return response.Result, nil
}

type CancelIdsResponse struct {
	SuccessIds   []string `json:"successIds"`
	FailIds      []string `json:"failIds"`
	TotalCount   int64    `json:"totalCount"`
	SuccessCount int64    `json:"successCount"`
	FailedCount  int64    `json:"failedCount"`
}

func (c *Connector) CancelByIds(symbol string, orderIds []string) (map[string]bool, error) {
	var response CancelIdsResponse
	var result = make(map[string]bool)
	err := c.Call(http.MethodPost, "spot/v4/cancel_orders", map[string]interface{}{
		"symbol":   symbol,
		"orderIds": orderIds,
	}, constants.Signed, &response)
	if err != nil {
		return nil, err
	}
	for _, id := range response.SuccessIds {
		result[id] = true
	}
	for _, id := range response.FailIds {
		result[id] = false
	}
	return result, nil
}
