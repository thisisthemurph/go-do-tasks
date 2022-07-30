package repository

import (
	"godo/internal/helper/ilog"
	"godo/internal/repository/entities"
)

type AccountQuery interface {
	GetAllAccounts() ([]*entities.Account, error)
	GetAccountById(accountId string) (*entities.Account, error)
	CreateAccount(newAccount *entities.Account) (*entities.Account, error)
	AccountExists(accountId string) (bool, error)
	AccountWithEmailAddressExists(email string) (bool, error)
}

type accountQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewAccountQuery(logger ilog.StdLogger) AccountQuery {
	return &accountQuery{log: logger}
}

func (q *accountQuery) GetAllAccounts() ([]*entities.Account, error) {
	q.log.Info("Fetching all accounts")

	var accounts []*entities.Account
	err := Database.Find(&accounts).Error
	ilog.ErrorlnIf(err, q.log)

	return accounts, err
}

func (q *accountQuery) GetAccountById(accountId string) (*entities.Account, error) {
	q.log.Infof("Fetching Account with ID %s", accountId)

	var account entities.Account
	err := Database.First(&account, "id = ?", accountId).Error
	ilog.ErrorlnIf(err, q.log)

	return &account, err
}

func (q *accountQuery) CreateAccount(newAccount *entities.Account) (*entities.Account, error) {
	q.log.Debugf("Creating Account with name of %s and email of %s", newAccount.Name)

	err := Database.Create(&newAccount).Error
	if err != nil {
		q.log.Error(err)
		return nil, err
	}

	return newAccount, err
}

func (q *accountQuery) AccountExists(accountId string) (bool, error) {
	var count int64
	r := Database.Model(&entities.Account{}).
		Where("id = ?", accountId).
		Count(&count)

	if r.Error != nil {
		q.log.Error(r.Error)
	}

	return count == 1, r.Error
}

func (q *accountQuery) AccountWithEmailAddressExists(email string) (bool, error) {
	var count int64
	r := Database.Model(&entities.Account{}).
		Where("email = ?", email).
		Count(&count)

	return count >= 1, r.Error
}
