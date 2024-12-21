package constants

import (
	"fmt"
	"regexp"
	"strings"
)

// StandardizeSymbol 统一的交易对格式：大写字母，下划线分隔 (例如：BTC_USDT)
func StandardizeSymbol(symbol string) (string, error) {
	// 移除所有非字母数字字符
	reg := regexp.MustCompile(`[^a-zA-Z0-9$]`)
	cleanedSymbol := reg.ReplaceAllString(symbol, "")

	// 将所有字母转换为大写
	cleanedSymbol = strings.ToUpper(cleanedSymbol)

	// 使用正则表达式匹配两个部分或多个部分
	reg = regexp.MustCompile(`([A-Z0-9]+)(USDT|USDC|BTC)`)
	matches := reg.FindStringSubmatch(cleanedSymbol)

	if len(matches) != 3 {
		return "", fmt.Errorf("invalid symbol format: %s", symbol)
	}
	return fmt.Sprintf("%s%s", matches[1], matches[2]), nil
}
