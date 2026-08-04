package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h[key]
	return value, exists
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))

	// no data found, assume we havent been given enough data and try again
	if idx == -1 {
		return 0, false, nil
	}

	line := string(data[:idx])
	consumed := idx + 2

	parts := strings.SplitN(line, ":", 2)

	if idx == 0 {
		return consumed, true, nil
	}
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header")
	}

	fieldName := parts[0]
	fieldValue := parts[1]

	fieldName = strings.TrimSpace(fieldName)
	fieldName = strings.ToLower(fieldName)

	fieldValue = strings.TrimSpace(fieldValue)

	if !isValidFieldName(fieldName) {
		return 0, false, errors.New("invalidg header field name")
	}

	if existingValue, exists := h[fieldName]; exists {
		h[fieldName] = existingValue + "," + fieldValue
	} else {
		h[fieldName] = fieldValue
	}

	return consumed, false, nil
}

func isValidFieldName(name string) bool {
	if len(name) == 0 {
		return false
	}

	for i := 0; i < len(name); i++ {
		if !isTokenChar(name[i]) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}

	if c >= 'A' && c <= 'Z' {
		return true
	}

	if c >= '0' && c <= '9' {
		return true
	}

	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+',
		'-', '.', '^', '_', '`', '|', '~':
		return true
	default:
		return false
	}
}
