package middleware

import (
	"context"
	"fmt"
	"godo/internal/api/httperror"
	"godo/internal/api/services"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	log         ilog.StdLogger
	authService services.AuthService
}

func NewAuthMiddleware(logger ilog.StdLogger, authService services.AuthService) AuthMiddleware {
	return AuthMiddleware{
		log:         logger,
		authService: authService,
	}
}

// Used to validate a request for a JWT
// The TokenRequest represents the auth data (e.g. emil and password) user to authenticate
func (m *AuthMiddleware) ValidateTokenRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tr auth.TokenRequest

		err := tr.FromJSON(r.Body)
		if err != nil {
			m.log.Error("The TokenRequest data was not in the expected JSON format")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate the project
		err = tr.Validate()
		if err != nil {
			m.log.Errorf("The TokenRequest failed validation: %s", err)
			e := httperror.New(
				http.StatusBadRequest,
				fmt.Sprintf("Error validating TokenRequest: %s", err),
			)

			http.Error(w, e.AsJson(), e.GetStatusCode())
			return
		}

		// Add the TokenRequest to the context
		ctx := context.WithValue(r.Context(), auth.TokenRequestKey{}, tr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Used to authenticate a given JWT
func (m *AuthMiddleware) AuthenticateRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedHttpError := httperror.New(http.StatusUnauthorized, "")
		tokenValue := r.Header.Get("Authorization")

		if tokenValue == "" {
			m.log.Info("No authorization token included in header")
			http.Error(w, "Unauthorized", unauthorizedHttpError.GetStatusCode())
			return
		}

		// Validate that the token string looks alright
		if !strings.HasPrefix(tokenValue, "Bearer") {
			m.log.Info("Bad token: the token does not start with `Bearer`")
			http.Error(w, "Unauthorized", unauthorizedHttpError.GetStatusCode())
			return
		}

		tokenParts := strings.SplitAfter(tokenValue, " ")
		if len(tokenParts) != 2 {
			m.log.Info("The bearer token does not follow the correct format")
			http.Error(w, "Unauthorized", unauthorizedHttpError.GetStatusCode())
			return
		}

		// Validate the auth token
		err := m.authService.ValidateToken(tokenParts[1])
		if err != nil {
			m.log.Info("The token is not valid")
			http.Error(w, "Unauthorized", unauthorizedHttpError.GetStatusCode())
			return
		}

		next.ServeHTTP(w, r)
	})
}
