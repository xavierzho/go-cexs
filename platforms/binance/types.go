package binance

type ErrorResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}
