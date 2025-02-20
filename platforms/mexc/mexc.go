package mexc

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.Mexc
}

func (c *Connector) SymbolPattern(symbol string) string {
	symbol, err := constants.StandardizeSymbol(symbol)
	if err != nil {
		return ""
	}
	return symbol
}

func NewConnector(cred *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: cred, Client: client}
}
