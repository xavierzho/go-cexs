package bitmart

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
	return constants.Bitmart
}

func NewConnector(base *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: base, Client: client}
}

func (c *Connector) SymbolPattern(symbol string) string {
	return constants.SymbolWithUnderline(symbol)
}
