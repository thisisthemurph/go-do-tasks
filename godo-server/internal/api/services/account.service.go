package services

import (
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type AccountService interface {
	CreateAccount(newAccount *entities.Account) (*entities.Account, error)
	AccountExists(accountId string) (bool, error)
	AccountWithEmailAddressExists(email string) (bool, error)
}

type accountService struct {
	log   ilog.StdLogger
	query repository.AccountQuery
}

func NewAccountService(accountQuery repository.AccountQuery, logger ilog.StdLogger) AccountService {
	return &accountService{
		query: accountQuery,
		log:   logger,
	}
}

func (a *accountService) CreateAccount(newAccount *entities.Account) (*entities.Account, error) {
	account, err := a.query.CreateAccount(newAccount)

	// TODO: Check if the account already exists

	if err != nil {
		a.log.Error("Issue creating Account: ", err)
		return nil, api.AccountNotCreatedError
	}

	return account, nil
}

func (a *accountService) AccountExists(accountId string) (bool, error) {
	return a.query.AccountExists(accountId)
}

func (a *accountService) AccountWithEmailAddressExists(email string) (bool, error) {
	return a.query.AccountWithEmailAddressExists(email)
}
