package handler

import (
	"encoding/json"
	"godo/internal/api/httperror"
	"godo/internal/api/httputils"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the login request
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Does the user exist?
	user, err := u.userService.GetUserByEmailAddress(request.Email)
	if err != nil {
		u.log.Warnf("User %s already exists", request.Email)
		http.Error(w, userAuthenticationError, http.StatusNotFound)
		return
	}

	// Verify the password
	err = user.VerifyPassword(request.Password)
	if err != nil {
		u.log.Warn("Bad authentication: incorrect password")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Gat a token for the user
	token, err := u.authService.GenerateJWT(user.Email, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponse{token}
	responseJson, _ := json.Marshal(response)
	httputils.MakeHttpResponse(w, string(responseJson), http.StatusOK)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the user
	user := auth.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	}

	err = u.userService.CreateUser(user)
	if err != nil {
		httpErr := httperror.New(http.StatusInternalServerError, "Unable to register user")
		http.Error(w, httpErr.AsJson(), httpErr.GetStatusCode())
	}

	// var userData string
	var d []byte
	d, _ = json.Marshal(user)

	httputils.MakeHttpResponse(w, string(d), http.StatusCreated)
}
