package utility

import (
	"fmt"
	"strings"
)

func StringIsIn(value string, slice []string) bool {
	for _, item := range slice {
		if strings.Compare(item, value) == 0 {
			return true
		}
	}
	return false
}

func PrettyPrintStringsSlice(slice []string) string {
	var result = "["
	for _, item := range slice {
		result += "\"" + fmt.Sprint(item) + "\"" + ","
	}
	return result[:len(result)-1] + "]"
}

func StringSliceEqual(slice []string, other []string) bool {
	if len(slice) != len(other) {
		return false
	}
	for index := range slice {
		if slice[index] != other[index] {
			return false
		}
	}
	return true
}
