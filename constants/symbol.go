package constants

import (
	"fmt"
	"regexp"
	"strings"
)

var UnifiedPattern = regexp.MustCompile(`([A-Z0-9]+)(USDT|BTC|USDC)`)

// StandardizeSymbol unifies the symbol format. For example: ETHUSDT
func StandardizeSymbol(symbol string) (string, error) {
	// Remove all non-alphanumeric characters
	reg := regexp.MustCompile(`[^a-zA-Z0-9]`)
	cleanedSymbol := reg.ReplaceAllString(symbol, "")

	// Convert all letters to uppercase
	cleanedSymbol = strings.ToUpper(cleanedSymbol)

	// Use a regular expression to match only two parts

	matches := UnifiedPattern.FindStringSubmatch(cleanedSymbol)
	if len(matches) != 3 {
		return "", fmt.Errorf("invalid symbol format: %s", symbol)
	}
	return fmt.Sprintf("%s%s", matches[1], matches[2]), nil
}

func SymbolWithUnderline(symbol string) string {

	matches := UnifiedPattern.FindStringSubmatch(symbol)
	if len(matches) != 3 {
		return ""
	}
	return fmt.Sprintf("%s_%s", matches[1], matches[2])
}

// SymbolWithHyphen returns the symbol with a hyphen separator. For example: ETH-USDT
func SymbolWithHyphen(symbol string) string {
	// First, standardize the symbol
	standardizedSymbol, err := StandardizeSymbol(symbol)
	if err != nil {
		return symbol // Return the original symbol if standardization fails
	}

	// Extract the two parts using the unified pattern
	matches := UnifiedPattern.FindStringSubmatch(standardizedSymbol)
	if len(matches) != 3 {
		return standardizedSymbol // Return the standardized symbol if extraction fails
	}

	// Combine the two parts with a hyphen
	return fmt.Sprintf("%s-%s", matches[1], matches[2])
}
