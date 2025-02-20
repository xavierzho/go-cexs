package utils

import (
	"fmt"
	"sort"
	"strings"
)

func EncodeParams(m map[string]any) string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result []string
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s=%v", key, m[key]))
	}
	return strings.Join(result, "&")
}
