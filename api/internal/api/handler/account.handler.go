package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	ehand "godo/internal/api/errorhandler"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type Accounts struct {
	log            ilog.StdLogger
	accountService services.AccountService
	userService    services.UserService
	eh             ehand.ErrorHandler
}

func NewAccountsHandler(
	logger ilog.StdLogger,
	accountService services.AccountService,
	userService services.UserService) Accounts {

	return Accounts{
		log:            logger,
		accountService: accountService,
		userService:    userService,
		eh:             ehand.New(),
	}
}

// CreateAccount Creates an account and appropriate user
func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	accountDto := r.Context().Value(entities.AccountKey{}).(dto.NewAccountDto)

	// Check that the account/user does not already exist
	// Neither shall be created whilst the other exists
	accountExists, err := a.accountService.AccountWithEmailAddressExists(accountDto.UserEmail)
	if err != nil {
		a.log.Error("Issue checking if account exists: ", err)
		api.ReturnError(ehand.ErrorAccountNotCreated, http.StatusInternalServerError, w)
		return
	}

	userExists, err := a.userService.UserWithEmailAddressExists(accountDto.UserEmail)
	if err != nil {
		a.log.Error("Issue checking if user exists: ", err)
		api.ReturnError(ehand.ErrorAccountNotCreated, http.StatusInternalServerError, w)
		return
	}

	if accountExists {
		api.ReturnError(ehand.ErrorAccountAlreadyExists, http.StatusBadRequest, w)
		return
	}

	if userExists {
		api.ReturnError(ehand.ErrorUserAlreadyExists, http.StatusFound, w)
		return
	}

	// Create the account
	newAccount := entities.Account{
		Name:  accountDto.Name,
		Email: accountDto.UserEmail,
	}

	createdAccount, err := a.accountService.CreateAccount(&newAccount)
	if status := a.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// Create the user
	newUser := entities.User{
		Name:      accountDto.UserName,
		Email:     accountDto.UserEmail,
		Username:  accountDto.UserUsername,
		Password:  accountDto.Password,
		AccountId: createdAccount.ID,
	}

	_, err = a.userService.CreateUser(newUser)
	if status := a.eh.HandleApiError(w, err); status != http.StatusOK {
		return
	}

	// TODO: Create a RespondWithPointer method to point to the newly created resource
	api.Respond(createdAccount, http.StatusCreated, w)
}
