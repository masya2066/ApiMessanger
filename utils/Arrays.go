package utils

import (
	"reflect"
	"strconv"
	"strings"
)

func IntSliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, val := range slice {
		strSlice[i] = strconv.FormatUint(uint64(val), 10)
	}
	return strings.Join(strSlice, ", ")
}

func IsArray(input []int) bool {
	valueType := reflect.TypeOf(input)
	if len(input) != 0 {
		if valueType.Kind() == reflect.Slice {
			elemType := valueType.Elem()
			return elemType.Kind() == reflect.Int
		}
	}
	return false
}
