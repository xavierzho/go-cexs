package constants

// Platform system identity exchange name
type Platform string

const (
	Binance Platform = "Binance"
	Bitmart Platform = "Bitmart"
	Mexc    Platform = "Mexc"
	LBank   Platform = "LBank"
	ByBit   Platform = "ByBit"
	Gate    Platform = "Gate"
	Okx     Platform = "Okx"
)

func (p Platform) String() string {
	return string(p)
}
