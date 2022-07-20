package httperror

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError interface {
	Error() string
	GetStatusCode() int
	GetStatusText() string
}

type httpError struct {
	Err          error  `json:"-"`
	ErrorMessage string `json:"errorMessage"`
	StatusCode   int    `json:"statusCode"`
	StatusText   string `json:"-"`
}

func New(status int, message string) HttpError {
	return &httpError{
		Err:          errors.New(message),
		ErrorMessage: message,
		StatusCode:   status,
		StatusText:   http.StatusText(status),
	}
}

func (e *httpError) Error() string {
	return fmt.Sprintf("Status %d: Err %v", e.StatusCode, e.Err)
}

func (e *httpError) GetStatusCode() int {
	return e.StatusCode
}

func (e *httpError) GetStatusText() string {
	return e.StatusText
}
