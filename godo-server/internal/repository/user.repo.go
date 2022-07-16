package repository

import (
	"errors"
	"fmt"
	"godo/internal/api"
	"godo/internal/auth"
	"godo/internal/helper/ilog"
)

type ApiUserQuery interface {
	CreateUser(user auth.User) (*auth.User, error)
	GetUserByEmailAddress(email string) (*auth.User, error)
	UserWithEmailAddressExists(email string) (bool, error)
}

type apiUserQuery struct {
	log ilog.StdLogger
}

func (d *dao) NewApiUserQuery(logger ilog.StdLogger) ApiUserQuery {
	return &apiUserQuery{log: logger}
}

func (q *apiUserQuery) CreateUser(newUser auth.User) (*auth.User, error) {
	q.log.Info("Registering new user")

	// Check if a user with the username or email already exists
	var foundUser auth.User
	result := Database.Where("email = ?", newUser.Email).Find(&foundUser)

	if result.RowsAffected >= 1 {
		q.log.Errorf("A user already exists with the given email address: %s", newUser.Email)
		return nil, api.UserNotFoundError
	}

	// Add the discriminator to the user
	discriminator := q.GetNextDiscriminator(newUser.Username)
	newUser.Discriminator = discriminator

	// Validate the user
	err := newUser.Validate()
	if err != nil {
		q.log.Warn("The user did not pass validation. ", err)
		return nil, errors.New(fmt.Sprintf("The user is not valid: %s", err.Error()))
	}

	// Hash the user's password
	if err := newUser.HashPassword(newUser.Password); err != nil {
		q.log.Error("Could not hash the user's password when creating user")
		return nil, err
	}

	// Insert the user into the database
	err = Database.Create(&newUser).Error
	ilog.ErrorlnIf(err, q.log)

	return &newUser, err
}

func (q *apiUserQuery) GetUserByEmailAddress(email string) (*auth.User, error) {
	var user auth.User
	err := Database.First(&user, "email = ?", email).Error

	if err != nil {
		q.log.Warn("There was an issue obtaining the user from the database")
		q.log.Error(err)
		return nil, api.UserNotFoundError
	}

	return &user, nil
}

func (q *apiUserQuery) UserWithEmailAddressExists(email string) (bool, error) {
	// _, err := q.GetUserByEmailAddress(email)

	var count int64
	r := Database.Model(&auth.User{}).
		Where("email = ?", email).
		Count(&count)

	return count >= 1, r.Error
}

func (q *apiUserQuery) GetNextDiscriminator(username string) uint32 {
	var result uint32
	row := Database.Table("users").
		Where("username = ?", username).
		Select("max(discriminator)").
		Row()

	err := row.Scan(&result)
	if err != nil {
		q.log.Infof("No user with username %s exists. Using discriminator of 1.", username)
		q.log.Error(err.Error())
		return 1
	}

	if result < 1 {
		q.log.Warnf("A less than 1 result of %d was returned for the discriminator", result)
		return 1
	}

	return result + 1
}
