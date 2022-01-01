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
		result += fmt.Sprint(item) + " ,"
	}
	return result[:len(result)-3] + "]"
}
