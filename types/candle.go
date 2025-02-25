package types

import (
	"strconv"
)

// CandleEntry fields verbose: [start_timestamp, open, high, low, close, volume, (volume_usd)]
type CandleEntry []float64

// CandlesEntry candle list
type CandlesEntry []CandleEntry

func Safe2Float(data any) float64 {
	switch v := data.(type) {
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case int:
		return float64(v)
	case uint64:
		return float64(v)
	case uint32:
		return float64(v)
	case uint:
		return float64(v)
	default:
		return 0
	}
}

type DepthEntry struct {
	Asks [][]string
	Bids [][]string
}
