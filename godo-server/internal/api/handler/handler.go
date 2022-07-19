package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"godo/internal/api"
	"godo/internal/repository/entities"
	"io"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// Converts a struct into a JSON object
func dataToJson(d interface{}) (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getParamFomRequest(r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	paramValue, paramExists := params[param]
	return paramValue, paramExists
}

func getStructFromContext[T any](ctx context.Context, key interface{}) T {
	return ctx.Value(key).(T)
}

func getUserFromContext(ctx context.Context) entities.User {
	return getStructFromContext[entities.User](ctx, entities.UserKey{})
}

type malformedRequest struct {
	status int
	err    error
}

func (mr *malformedRequest) Error() string {
	return mr.err.Error()
}

// Attempts to decode the JSON body, returning a custom malformed request error type.
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, err: errors.New(msg)}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, err: errors.New(msg)}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}
	}

	return nil
}

func handleMalformedJSONError(w http.ResponseWriter, err error) {
	var mr *malformedRequest
	if errors.As(err, &mr) {
		api.ReturnError(mr.err, mr.status, w)
	} else {
		api.ReturnError(errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError, w)
	}
}
