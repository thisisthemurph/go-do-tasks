package repository

import (
	"errors"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
)

type ApiUserQuery interface {
	CreateUser(user auth.User) error
	GetUserByEmailAddress(email string) (auth.User, error)
	UserWithEmailAddressExists(email string) (bool, error)
}

type apiUserQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewApiUserQuery(logger ilog.StdLogger) ApiUserQuery {
	return &apiUserQuery{log: logger}
}

func (q *apiUserQuery) CreateUser(newUser auth.User) error {
	q.log.Info("Registering new user")

	// Check if a user with the username or email already exists
	var foundUser auth.User
	result := Database.Where("email = ?", newUser.Email).Find(&foundUser)

	if result.RowsAffected >= 1 {
		q.log.Errorf("A user already exists with the given email address: %s", newUser.Email)
		return errors.New("A user with the given email address already exists")
	}

	// TODO: Validate the user struct

	// Insert the user into the database
	if err := newUser.HashPassword(newUser.Password); err != nil {
		q.log.Error("Could not hash the user's password when creating user")
		return err
	}

	err := Database.Create(&newUser).Error
	ilog.ErrorlnIf(err, q.log)

	return err
}

func (q *apiUserQuery) GetUserByEmailAddress(email string) (user auth.User, err error) {
	err = Database.First(&user, "email = ?", email).Error

	if err != nil {
		q.log.Error("There was an issue obtaining the user from the database")
		q.log.Error(err)
	}

	return
}

func (q *apiUserQuery) UserWithEmailAddressExists(email string) (bool, error) {
	_, err := q.GetUserByEmailAddress(email)
	return err != nil, err
}
