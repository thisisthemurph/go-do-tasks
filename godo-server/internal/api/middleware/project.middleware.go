package middleware

import (
	"context"
	"fmt"
	"godo/internal/api/httperror"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type ProjectMiddleware struct {
	log ilog.StdLogger
}

func NewProjectMiddleware(logger ilog.StdLogger) ProjectMiddleware {
	return ProjectMiddleware{log: logger}
}

func (m *ProjectMiddleware) ValidateProjectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var project entities.Project

		err := project.FromJSON(r.Body)
		if err != nil {
			m.log.Error("The Project data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the project
		err = project.Validate()
		if err != nil {
			m.log.Errorf("The Project failed validation: %s", err)
			e := httperror.New(
				http.StatusBadRequest,
				fmt.Sprintf("Error validating project: %s", err),
			)

			http.Error(w, e.AsJson(), e.GetStatusCode())
			return
		}

		// Add the project to the context
		ctx := context.WithValue(r.Context(), entities.ProjectKey{}, project)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
