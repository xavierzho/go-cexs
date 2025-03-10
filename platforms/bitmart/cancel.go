package bitmart

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
)

func (c *Connector) CancelAll(symbol string) error {
	var response map[string]interface{}
	err := c.Call(http.MethodPost, CancelAllEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
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
	err := c.Call(http.MethodPost, CancelEndpoint, &platforms.ObjectBody{
		SymbolFiled: symbol,
		"order_id":  orderId,
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
	err := c.Call(http.MethodPost, CancelsEndpoint, &platforms.ObjectBody{
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
