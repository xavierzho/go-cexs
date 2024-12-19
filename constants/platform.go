package constants

// Platform exchange name
type Platform string

const (
	Binance Platform = "Binance"
	Bitmart Platform = "Bitmart"
	CoinW   Platform = "CoinW"
	Mexc    Platform = "Mexc"
	LBank   Platform = "LBank"
	XT      Platform = "XT"
	ByBit   Platform = "ByBit"
)

func (p Platform) String() string {
	return string(p)
}
