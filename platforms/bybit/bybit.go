package bybit

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"
	"strconv"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func (c *Connector) Name() constants.Platform {
	return constants.ByBit
}

func (c *Connector) SymbolPattern(symbol string) string {
	symbol, _ = constants.StandardizeSymbol(symbol)
	return symbol
}

func NewConnector(cred *platforms.Credentials, client *http.Client) platforms.SpotConnector {
	return &Connector{Credentials: cred, Client: client}
}

func timeConvert(interval string) string {
	unit := interval[len(interval)-1]
	value, _ := strconv.ParseInt(interval[:len(interval)-1], 10, 64)
	switch string(unit) {
	case "m":
		return strconv.FormatInt(value, 10)
	case "h":
		return strconv.FormatInt(value*60, 10)
	case "M":
		return "M"
	case "D":
		return "D"
	case "W":
		return "W"
	default:
		return "1"
	}
}
