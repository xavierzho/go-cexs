package types

// CandleBaseKeys for CandleEntry.FromList keys params
var CandleBaseKeys = []string{
	"time_start", "open", "high", "low", "close", "volume",
}

//	type CandleEntry struct {
//		O  string // open
//		H  string // high
//		L  string // low
//		C  string // close
//		V  string // volume
//		Ts int64  // timestamp start
//		Te int64  // timestamp end
//	}
type CandleEntry map[string]any

func (c *CandleEntry) FromList(data []any, keys []string) {
	var m = make(map[string]any)
	if len(data) != len(keys) {
		return
	}
	for i, key := range keys {
		m[key] = data[i]
	}
	*c = m
}

type DepthEntry struct {
	Asks [][]string
	Bids [][]string
}
