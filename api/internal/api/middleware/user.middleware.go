package middleware

import (
	"context"
	"fmt"
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/helper/validate"
	"godo/internal/repository/entities"
	"net/http"
)

type UserMiddleware struct {
	log ilog.StdLogger
}

func NewUserMiddleware(logger ilog.StdLogger) ProjectMiddleware {
	return ProjectMiddleware{log: logger}
}

func (m *UserMiddleware) ValidateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user entities.User

		err := api.FromJSON(user, r.Body)
		if err != nil {
			m.log.Error("The User data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the project
		err = validate.Struct(user)
		if err != nil {
			m.log.Errorf("The User failed validation: %s", err)

			e := fmt.Errorf("error validating user: %s", err)
			api.ReturnError(e, http.StatusBadRequest, w)
			return
		}

		// Add the project to the context
		ctx := context.WithValue(r.Context(), entities.ProjectKey{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
