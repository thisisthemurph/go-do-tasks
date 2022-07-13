package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"godo/internal/api/dto"
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

func (m *ProjectMiddleware) ValidateNewProjectDtoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var projectDto dto.NewProjectDto

		err := json.NewDecoder(r.Body).Decode(&projectDto)
		if err != nil {
			m.log.Error("The Project data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = projectDto.Validate()
		if err != nil {
			m.log.Errorf("The Project failed validation: %s", err)
			e := httperror.New(
				http.StatusBadRequest,
				fmt.Sprintf("Error validating project: %s", err),
			)

			http.Error(w, e.AsJson(), e.GetStatusCode())
			return
		}

		// Add the project dto to the context
		m.log.Info("Adding NewProjectDto to context ", projectDto)
		ctx := context.WithValue(r.Context(), entities.ProjectKey{}, projectDto)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
