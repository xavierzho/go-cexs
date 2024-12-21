package bitmart

import (
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"regexp"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.Bitmart
}

func NewConnector(base *platforms.Credentials, client *http.Client) platforms.Connector {
	return &Connector{Credentials: base, Client: client}
}

func (c *Connector) SymbolPattern(symbol string) string {
	reg := regexp.MustCompile(`([A-Z0-9])(USDT|BTC|USDC)`)
	matches := reg.FindStringSubmatch(symbol)
	return fmt.Sprintf("%s_%s", matches[1], matches[2])
}
