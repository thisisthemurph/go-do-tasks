package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/helper/ilog"
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
			api.ReturnError(err, http.StatusBadRequest, w)
			return
		}

		err = projectDto.Validate()
		if err != nil {
			m.log.Errorf("The Project failed validation: %s", err)

			e := errors.New(fmt.Sprintf("Error validating project: %s", err.Error()))
			api.ReturnError(e, http.StatusBadRequest, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
