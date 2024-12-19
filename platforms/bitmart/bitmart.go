package bitmart

import (
	cexconns "github.com/xavierzho/go-cexs"
	"github.com/xavierzho/go-cexs/constants"
	"net/http"
)

type Connector struct {
	*cexconns.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.Bitmart
}

func NewConnector(base *cexconns.Credentials, client *http.Client) cexconns.Connector {
	return &Connector{Credentials: base, Client: client}
}
