package cexconns

import (
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/platforms/binance"
	"github.com/xavierzho/go-cexs/platforms/bitmart"
	"net/http"
)

func NewExchange(ex constants.Platform, apikey, apiSecret string, option *string) platforms.Connector {
	switch ex {
	case constants.Bitmart:
		return bitmart.NewConnector(platforms.NewCredentials(apikey, apiSecret, option), &http.Client{})
	case constants.Binance:
		return binance.NewConnector(platforms.NewCredentials(apikey, apiSecret, option), &http.Client{})
	default:
		return nil
	}
}

func NewMarketStream(ex constants.Platform) platforms.MarketStreamer {
	switch ex {
	case constants.Binance:
		return binance.NewMarketStream()
	case constants.Bitmart:
		return bitmart.NewMarketStream()
	default:
		return nil
	}
}

func NewUserDataStream(ex constants.Platform, apikey, apiSecret string, option *string) platforms.UserDataStreamer {
	switch ex {
	case constants.Binance:
		return binance.NewUserStream(apikey, apiSecret)
	case constants.Bitmart:
		return bitmart.NewUserStream(platforms.NewCredentials(apikey, apiSecret, option))
	default:
		return nil
	}
}
