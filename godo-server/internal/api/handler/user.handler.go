package handler

import (
	"encoding/json"
	"godo/internal/api"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"

	"github.com/go-playground/validator"
)

type Users struct {
	log            ilog.StdLogger
	authService    services.AuthService
	accountService services.AccountService
	userService    services.UserService
}

func NewUsersHandler(
	logger ilog.StdLogger,
	authService services.AuthService,
	accountService services.AccountService,
	userService services.UserService) *Users {

	return &Users{log: logger, authService: authService, accountService: accountService, userService: userService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Validate the login request
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Does the user exist?
	user, err := u.userService.GetUserByEmailAddress(request.Email)
	switch err {
	case nil:
		break
	case api.UserNotFoundError:
		api.ReturnError(api.UserAuthenticationError, http.StatusUnauthorized, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	// Verify the password
	err = user.VerifyPassword(request.Password)
	if err != nil {
		u.log.Warn("Bad authentication: incorrect password")
		api.ReturnError(api.UserAuthenticationError, http.StatusUnauthorized, w)
		return
	}

	// Gat a token for the user
	token, err := u.authService.GenerateJWT(user.Email, user.Username, user.AccountId)
	if err != nil {
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	response := LoginResponse{token}
	api.Respond(response, http.StatusOK, w)
}

type RegistrationRequest struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	AccountId string `json:"account_id" validate:"required"`
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	var request RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Validate the registration body data
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Ensure the account exists
	accountExists, err := u.accountService.AccountExists(request.AccountId)
	if !accountExists {
		api.ReturnError(api.AccountNotFoundError, http.StatusNotFound, w)
		return
	}

	// Create the newUser
	newUser := entities.User{
		Name:      request.Name,
		Email:     request.Email,
		Username:  request.Username,
		Password:  request.Password,
		AccountId: request.AccountId,
	}

	var createdUser *entities.User
	createdUser, err = u.userService.CreateUser(newUser)

	switch err {
	case nil:
		break
	case api.UserAlreadyExistsError:
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(createdUser, http.StatusOK, w)
}
