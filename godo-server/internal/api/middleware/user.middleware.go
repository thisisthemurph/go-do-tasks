package middleware

import (
	"context"
	"fmt"
	"godo/internal/api/httperror"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type UserMiddleware struct {
	log ilog.StdLogger
}

func NewUserMiddleware(logger ilog.StdLogger) ProjectMiddleware {
	return ProjectMiddleware{log: logger}
}

func (m *ProjectMiddleware) ValidateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user auth.User

		err := user.FromJSON(r.Body)
		if err != nil {
			m.log.Error("The User data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the project
		err = user.Validate()
		if err != nil {
			m.log.Errorf("The User failed validation: %s", err)
			e := httperror.New(
				http.StatusBadRequest,
				fmt.Sprintf("Error validating project: %s", err),
			)

			http.Error(w, e.AsJson(), e.GetStatusCode())
			return
		}

		// Add the project to the context
		ctx := context.WithValue(r.Context(), entities.ProjectKey{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
