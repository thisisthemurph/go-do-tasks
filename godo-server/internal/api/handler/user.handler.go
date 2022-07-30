package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/helper/validate"
	"godo/internal/repository/entities"
	"net/http"
)

type Users struct {
	log            ilog.StdLogger
	authService    services.AuthService
	accountService services.AccountService
	userService    services.UserService
	eh             ehand.ErrorHandler
}

func NewUsersHandler(
	logger ilog.StdLogger,
	authService services.AuthService,
	accountService services.AccountService,
	userService services.UserService) *Users {

	return &Users{
		log:            logger,
		authService:    authService,
		accountService: accountService,
		userService:    userService,
		eh:             ehand.New(),
	}
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	request, err := getDtoFromJSONBody[dto.LoginRequestDto](w, r)
	if err != nil {
		return
	}

	// Validate the login request
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Does the user exist?
	user, err := u.userService.GetUserByEmailAddress(request.Email)
	if status := u.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Verify the password
	err = user.VerifyPassword(request.Password)
	if err != nil {
		u.log.Debug("Bad authentication: incorrect password")
		api.ReturnError(ehand.ErrorUserAuthentication, http.StatusUnauthorized, w)
		return
	}

	// Gat a token for the user
	token, err := u.authService.GenerateJWT(user.Email, user.Username, user.AccountId)
	if status := u.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	response := dto.LoginResponseDto{Token: token}
	api.Respond(response, http.StatusOK, w)
}

func (u *Users) Register(w http.ResponseWriter, r *http.Request) {
	request, err := getDtoFromJSONBody[dto.RegistrationRequestDto](w, r)
	if err != nil {
		return
	}

	// Validate the registration body data
	err = validate.Struct(request)
	if err != nil {
		api.ReturnError(err, http.StatusBadRequest, w)
		return
	}

	// Ensure the account exists
	accountExists, err := u.accountService.AccountExists(request.AccountId)
	if status := u.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	if !accountExists {
		api.ReturnError(ehand.ErrorAccountNotFound, http.StatusNotFound, w)
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
	if status := u.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	api.Respond(createdUser, http.StatusOK, w)
}
