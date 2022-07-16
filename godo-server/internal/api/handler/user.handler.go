package handler

import (
	"encoding/json"
	"godo/internal/api"
	"godo/internal/api/services"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"net/http"

	"github.com/go-playground/validator"
)

type Users struct {
	log         ilog.StdLogger
	authService services.AuthService
	userService services.UserService
}

func NewUsersHandler(logger ilog.StdLogger, authService services.AuthService, userService services.UserService) *Users {
	return &Users{log: logger, authService: authService, userService: userService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

var userAuthenticationError string = "A user with the given username and password combination could not be found"

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ReturnError(err.Error(), http.StatusBadRequest, w)
		return
	}

	// Validate the login request
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err.Error(), http.StatusBadRequest, w)
		return
	}

	// Does the user exist?
	user, err := u.userService.GetUserByEmailAddress(request.Email)
	switch err {
	case nil:
		break
	case api.UserNotFoundError:
		api.ReturnError(err.Error(), http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err.Error(), http.StatusInternalServerError, w)
		return
	}

	// Verify the password
	err = user.VerifyPassword(request.Password)
	if err != nil {
		u.log.Warn("Bad authentication: incorrect password")
		api.ReturnError(err.Error(), http.StatusUnauthorized, w)
		return
	}

	// Gat a token for the user
	token, err := u.authService.GenerateJWT(user.Email, user.Username)
	if err != nil {
		api.ReturnError(err.Error(), http.StatusInternalServerError, w)
		return
	}

	response := LoginResponse{token}
	api.Respond(response, http.StatusOK, w)
}

type RegistrationRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	var request RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ReturnError(err.Error(), http.StatusBadRequest, w)
		return
	}

	// Validate the registration body data
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err.Error(), http.StatusBadRequest, w)
		return
	}

	// Create the newUser
	newUser := auth.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	}

	var createdUser *auth.User
	createdUser, err = u.userService.CreateUser(newUser)

	switch err {
	case nil:
		break
	case api.UserAlreadyExistsError:
		api.ReturnError(err.Error(), http.StatusNotFound, w)
		return
	default:
		api.ReturnError(err.Error(), http.StatusInternalServerError, w)
		return
	}

	api.Respond(createdUser, http.StatusOK, w)
}
