package handler

import (
	"godo/internal/api"
	"godo/internal/api/dto"
	"godo/internal/api/services"
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
	"net/http"
)

type Accounts struct {
	log            ilog.StdLogger
	accountService services.AccountService
	userService    services.UserService
}

func NewAccountsHandler(
	logger ilog.StdLogger,
	accountService services.AccountService,
	userService services.UserService) Accounts {

	return Accounts{
		log:            logger,
		accountService: accountService,
		userService:    userService,
	}
}

func (a *Accounts) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var accountDto dto.NewAccountDto
	accountDto = r.Context().Value(entities.AccountKey{}).(dto.NewAccountDto)

	// TODO: Check if the account already exists
	// TODO: Check if the user already exists

	newAccount := entities.Account{
		Name:  accountDto.Name,
		Email: accountDto.UserEmail,
	}

	createdAccount, err := a.accountService.CreateAccount(&newAccount)
	if err != nil {
		api.ReturnError(api.AccountNotCreatedError, http.StatusInternalServerError, w)
		return
	}

	newUser := entities.User{
		Name:      accountDto.UserName,
		Email:     accountDto.UserEmail,
		Username:  accountDto.UserUsername,
		Password:  accountDto.Password,
		AccountId: createdAccount.ID,
	}

	_, err = a.userService.CreateUser(newUser)
	switch err {
	case nil:
		break
	case api.UserAlreadyExistsError:
		api.ReturnError(err, http.StatusFound, w)
		return
	default:
		api.ReturnError(err, http.StatusInternalServerError, w)
		return
	}

	api.Respond(createdAccount, http.StatusOK, w)
}
