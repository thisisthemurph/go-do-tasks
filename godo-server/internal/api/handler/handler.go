package handler

import (
	"encoding/json"
)

// Converts a struct into a JSON object
func dataToJson(d interface{}) (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
