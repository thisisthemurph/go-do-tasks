package middleware

import (
	"log"
	"net/http"
)

type GenericMiddleware struct{}

func (m *GenericMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: %v\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
