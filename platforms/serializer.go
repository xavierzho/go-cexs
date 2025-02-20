package platforms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Serializer interface {
	Serialize() (io.Reader, error)
	EncodeQuery() (string, error)
	Set(string, any)
	Exists(string) (any, bool)
}

type ObjectBody map[string]any

func (m *ObjectBody) Serialize() (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bs), nil
}
func (m *ObjectBody) EncodeQuery() (string, error) {
	var keys []string
	for key := range *m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result []string
	for _, key := range keys {
		result = append(result, fmt.Sprintf("%s=%v", key, (*m)[key]))
	}
	return strings.Join(result, "&"), nil
}
func (m *ObjectBody) Set(key string, val any) {
	(*m)[key] = val
}
func (m *ObjectBody) Exists(key string) (any, bool) {
	v, ok := (*m)[key]
	return v, ok
}

type ArrayBody []map[string]any

func (m *ArrayBody) Serialize() (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bs), nil
}

func (m *ArrayBody) EncodeQuery() (string, error) {
	return "", fmt.Errorf("not support this method")
}

func (m *ArrayBody) Set(key string, val any) {
	idx, err := strconv.Atoi(key)
	if err != nil {
		return
	}
	v, ok := val.(map[string]any)
	if !ok {
		return
	}
	if idx >= len(*m) {
		*m = append(*m, v)
	} else {
		(*m)[idx] = v
	}

}
func (m *ArrayBody) Exists(key string) (any, bool) {
	idx, err := strconv.Atoi(key)
	if err != nil {
		return nil, false
	}
	if idx < 0 || idx >= len(*m) {
		return nil, false
	}
	return (*m)[idx], true
}
