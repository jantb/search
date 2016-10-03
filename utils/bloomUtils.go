package utils

import (
	"strings"
	"unicode"
)

func GetBloomKeysFromLine(line string)[][]byte{
	fields := strings.FieldsFunc(line, func(r rune) bool{
		if r == '=' {
			return false
		}
		return unicode.IsSpace(r) || unicode.IsSymbol(r) || unicode.IsPunct(r)
	})
	keys := [][]byte{}
	for _, field := range fields {
		keys = append(keys, []byte(field))
	}
	return keys
}
