package handler

import (
	"fmt"
	"godo/internal/api/httputils"
	"godo/internal/api/services"
	"godo/internal/auth"
	"log"
	"net/http"
)

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	var status int = http.StatusOK

	h.log.Infof("AuthHandler. Method: %v", r.Method)

	switch r.Method {
	case http.MethodPost:
		result, status = generateToken(h.userService, h.authService, r)
		break
	default:
		h.log.Errorf("Bad method %v", r.Method)
		result = fmt.Sprintf("Bad method %s", r.Method)
		status = http.StatusBadGateway
		break
	}

	httputils.MakeHttpResponse(w, result, status)
}

func generateToken(userService services.UserService, authService services.AuthService, r *http.Request) (result string, status int) {
	var request auth.TokenRequest
	request.FromHttpRequest(r)

	result = "Username or password combination is incorrect"
	status = http.StatusUnauthorized

	// Get the user and verify they exist
	user, err := userService.GetUserByEmailAddress(request.Email)
	if err != nil {
		log.Printf("Could not find user with email address %v\n", request.Email)
		return
	}

	// Verify the user's password is correct
	err = user.VerifyPassword(request.Password)
	if err != nil {
		log.Printf("Could not verify the given password")
		return
	}

	// Generate a token for the user
	token, err := authService.GenerateJWT(request.Email, "")
	if err != nil {
		log.Printf("Could not generate token")
		return
	}

	result = token
	status = http.StatusOK
	return
}
