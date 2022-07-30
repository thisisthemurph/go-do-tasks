package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"godo/internal/api"
	"godo/internal/helper/validate"
	"godo/internal/repository/entities"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// Fetches the given parameter from the HTTP request object
// Returns a string value and a success boolean
func getParamFomRequest(r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	value, exists := params[param]
	if !exists {
		log.Printf("the param %s does not exist in the request\n", param)
	}

	return value, exists
}

func getUintParamFomRequest(r *http.Request, param string) (uint, bool) {
	value, exists := getParamFomRequest(r, param)
	if !exists {
		return 0, false
	}

	p, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		log.Printf("Error processing param %s: %s", param, err.Error())
		return 0, false
	}

	return uint(p), exists
}

func getStructFromContext[T any](ctx context.Context, key interface{}) T {
	return ctx.Value(key).(T)
}

// Fetches the user from the given HTTP Request context
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
			errMsg := fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, err: errMsg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			errMsg := fmt.Errorf("request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, err: errMsg}

		case errors.As(err, &unmarshalTypeError):
			errMsg := fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, err: errMsg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			errMsg := fmt.Errorf("request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, err: errMsg}

		case errors.Is(err, io.EOF):
			errMsg := errors.New("request body must not be empty")
			return &malformedRequest{status: http.StatusBadRequest, err: errMsg}

		case err.Error() == "http: request body too large":
			errMsg := errors.New("request body must not be larger than 1MB")
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, err: errMsg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, err: errors.New(msg)}
	}

	return nil
}

func handleMalformedJSONError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var mr *malformedRequest
	if errors.As(err, &mr) {
		api.ReturnError(mr.err, mr.status, w)
	} else {
		api.ReturnError(errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError, w)
	}
}

func getDtoFromJSONBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var obj T
	err := decodeJSONBody(w, r, &obj)
	if err != nil {
		handleMalformedJSONError(w, err)
		return nil, err
	}

	err = validate.Struct(obj)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return nil, err
	}

	return &obj, nil
}
