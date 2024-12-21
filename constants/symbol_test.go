package constants

import (
	"fmt"
	"testing"
)

func TestStandardizeSymbol(t *testing.T) {
	testSymbols := []string{
		"BTCUSDT",
		"BTC_USDT",
		"btc_usdt",
		"ethusdt",
		"ETH/USDT",
		"btc-usdt",
		"XRP/BTC",
		"BNB_BTC",
		"invalid-symbol", // 测试无效符号
		"BTC/USD/TEST",   //测试多斜杠符号
	}

	for _, symbol := range testSymbols {
		standardized, err := StandardizeSymbol(symbol)
		if err != nil {
			fmt.Printf("Symbol: %s, Error: %v\n", symbol, err)
		} else {
			fmt.Printf("Symbol: %s, Standardized: %s\n", symbol, standardized)
		}
	}
}
