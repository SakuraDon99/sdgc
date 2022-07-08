package datahelper

import "encoding/json"

func Marshal(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return "Marshal failed"
	}
	return string(b)
}

func Unmarshal[T any](jsonStr string) (T, error) {
	var tar T
	err := json.Unmarshal([]byte(jsonStr), &tar)
	return tar, err
}
