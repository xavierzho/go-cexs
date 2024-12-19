package constants

import "fmt"

type SymbolPattern string

const (
	Underscore  SymbolPattern = "%s_%s"
	NoSeparator SymbolPattern = "%s%s"
)

// Format if you need before build request struct using strings.ToLower or strings.ToUpper to convert
func (sp SymbolPattern) Format(base, quote string) string {
	return fmt.Sprintf(string(sp), base, quote)
}

var SymbolPatternMap = map[Platform]SymbolPattern{
	Bitmart: Underscore,
}
