package gate

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
	return constants.Gate
}

func (c *Connector) SymbolPattern(symbol string) string {
	return constants.SymbolWithUnderline(symbol)
}

func NewConnector(cred *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: cred, Client: client}
}
