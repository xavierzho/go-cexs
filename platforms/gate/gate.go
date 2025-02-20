package gate

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.Gate
}

func (c *Connector) SymbolPattern(symbol string) string {
	reg := regexp.MustCompile(`([A-Z0-9])(USDT|BTC|USDC)`)
	matches := reg.FindStringSubmatch(symbol)
	if len(matches) != 3 {
		return ""
	}
	return fmt.Sprintf("%s_%s", matches[1], matches[2])
}

func NewConnector(cred *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: cred, Client: client}
}
