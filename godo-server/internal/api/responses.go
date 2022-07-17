package api

import (
	"godo/internal/api/httperror"
	"net/http"
)

func Respond(i interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	ToJSON(i, w)
}

func ReturnError(message error, status int, w http.ResponseWriter) {
	httpErr := httperror.New(status, message.Error())
	Respond(httpErr, status, w)
}
