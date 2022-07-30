package errorhandler

import (
	"fmt"
	"godo/internal/api"
	"net/http"
)

type ErrorHandler struct {
	em map[error]int
}

func New() ErrorHandler {
	return ErrorHandler{
		em: makeErrorMap(),
	}
}

func (e ErrorHandler) GetStatus(err error) (int, error) {
	status, prs := e.em[err]
	if !prs {
		return -1, fmt.Errorf("the given error does not exist in the error map")
	}

	return status, nil
}

func (e ErrorHandler) HandleApiError(w http.ResponseWriter, err error) int {
	if err == nil {
		return http.StatusOK
	}

	status, _ := e.GetStatus(err)
	api.ReturnError(err, status, w)

	return status
}
