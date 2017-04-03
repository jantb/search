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
		f :=[]rune{}
		for _, r := range field {
			f = append(f, r)
			keys = append(keys, []byte(string(f)))
		}
	}
	return keys
}
