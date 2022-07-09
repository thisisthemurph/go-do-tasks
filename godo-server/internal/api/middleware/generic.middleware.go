package middleware

import (
	"godo/internal/helper/ilog"
	"net/http"
)

type GenericMiddleware struct {
	log ilog.StdLogger
}

func NewGenericMiddleware(logger ilog.StdLogger) GenericMiddleware {
	return GenericMiddleware{log: logger}
}

func (m *GenericMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Infof("%v: %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
