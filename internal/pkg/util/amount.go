package util

import "strings"

// IsDecimalValid checks the num is positive value and with up to two decimal places.
func IsDecimalValid(num string) bool {
	numArr := strings.Split(num, ".")
	if len(numArr) == 1 {
		return true
	}
	if len(numArr) == 2 && len(numArr[1]) <= 2 {
		return true
	}
	return false
}
