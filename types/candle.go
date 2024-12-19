package types

//	type Candle struct {
//		O  string // open
//		H  string // high
//		L  string // low
//		C  string // close
//		V  string // volume
//		Ts int64  // timestamp start
//		Te int64  // timestamp end
//	}
type Candle map[string]any

func (c *Candle) FromList(data []any, keys []string) {
	var m = make(map[string]any)
	if len(data) != len(keys) {
		return
	}
	for i, key := range keys {
		m[key] = data[i]
	}
	*c = m
}
