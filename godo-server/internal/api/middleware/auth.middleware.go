package middleware

import (
	"context"
	"errors"
	"fmt"
	"godo/internal/api/services"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type AuthMiddleware struct {
	log         ilog.StdLogger
	authService services.AuthService
	userService services.UserService
}

func NewAuthMiddleware(
	logger ilog.StdLogger,
	authService services.AuthService,
	userService services.UserService) AuthMiddleware {
	return AuthMiddleware{
		log:         logger,
		authService: authService,
		userService: userService,
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

			e := errors.New(fmt.Sprintf("Error validating TokenRequest: %s", err))
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		// Add the TokenRequest to the context
		m.log.Info("Adding auth TokenRequest to context ", tr)
		ctx := context.WithValue(r.Context(), auth.TokenRequestKey{}, tr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Used to authenticate a given JWT
func (m *AuthMiddleware) AuthenticateRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValue := r.Header.Get("Authorization")

		// Convert the token from `Bearer abc123` to `abc123`
		// This also validates the given token string value
		token, err := m.authService.BearerTokenToToken(tokenValue)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the auth token
		err = m.authService.ValidateTokenClaims(token)
		if err != nil {
			m.log.Info("The token is not valid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch the appropriate user
		claims, err := m.authService.GetClaims(token)
		if err != nil {
			m.log.Error("Could not get claims from the token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.userService.GetUserByEmailAddress(claims.Email)
		if err != nil {
			m.log.Error("Could not determine user with email address %s from signed token", claims.Email)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify the user's JWT account_id is the same as in the database
		if user.AccountId != claims.AccountId {
			m.log.Warnf("The user account in the token claims is incorrect. Expected %s but got %s", user.AccountId, claims.AccountId)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach the user to the context
		m.log.Info("Adding user to context ", user)
		ctx := context.WithValue(r.Context(), entities.UserKey{}, *user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
