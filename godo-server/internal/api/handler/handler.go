package handler

import (
	"context"
	"encoding/json"
	"godo/internal/repository/entities"
	"net/http"

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
	// return ctx.Value(entities.UserKey{}).(entities.User)
	return getStructFromContext[entities.User](ctx, entities.UserKey{})
}
