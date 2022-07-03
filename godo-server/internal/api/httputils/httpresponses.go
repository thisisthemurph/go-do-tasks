package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeHttpResponse(w http.ResponseWriter, result string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, result)
}

type HttpError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"errorMessage"`
}

func (h HttpError) String() (string) {
	data, _ := json.Marshal(h)
	return string(data)
}

func MakeHttpError(statusCode int, message string) (string, int) {
	httpErr := HttpError{statusCode, message}
	return httpErr.String(), statusCode
}