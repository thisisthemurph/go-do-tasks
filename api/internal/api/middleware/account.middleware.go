package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type AccountMiddleware struct {
	log ilog.StdLogger
}

func NewAccountMiddleware(logger ilog.StdLogger) AccountMiddleware {
	return AccountMiddleware{log: logger}
}

func (m *AccountMiddleware) ValidateNewAccountDtoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accountDto dto.NewAccountDto

		err := json.NewDecoder(r.Body).Decode(&accountDto)
		if err != nil {
			m.log.Error("The Account data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = accountDto.Validate()
		if err != nil {
			m.log.Errorf("The Account failed validation: %s", err)

			validationErr := errors.New(fmt.Sprintf("Error validating Account: %s", err.Error()))
			api.ReturnError(validationErr, http.StatusBadRequest, w)
			return
		}

		// Add the Account dto to the context
		m.log.Info("Adding NewAccountDto to context ", accountDto)
		ctx := context.WithValue(r.Context(), entities.AccountKey{}, accountDto)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
