package datahelper

import "strconv"

func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func StartWith(str string, sub string) bool {
	if len(str) < len(sub) {
		return false
	}
	return str[:len(sub)] == sub
}

func StringIn(str string, arr ...string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
