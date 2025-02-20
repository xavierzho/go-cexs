package okx

import (
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.Okx
}

func (c *Connector) SymbolPattern(symbol string) string {
	return constants.SymbolWithHyphen(symbol)
}

func NewConnector(cred *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: cred, Client: client}
}

type RestReturn[T fmt.Stringer] struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	InTime  string `json:"inTime"`
	OutTime string `json:"outTime"`
	Data    []T    `json:"data"`
}
