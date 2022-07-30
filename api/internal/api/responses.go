package api

import (
	"godo/internal/api/httperror"
	"log"
	"net/http"
)

func Respond(i interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := ToJSON(i, w)
	if err != nil {
		log.Println("Could not format JSON response: ", err.Error())
	}
}

func ReturnError(message error, status int, w http.ResponseWriter) {
	httpErr := httperror.New(status, message.Error())
	Respond(httpErr, status, w)
}
