package utils

import "encoding/json"

func InterfaceToString(input interface{}) string {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return ""
	}
	return string(jsonInput)
}
