package services

import (
	"godo/internal/auth"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
)

type UserService interface {
	GetUserByEmailAddress(email string) (user auth.User, err error)
	UserWithEmailAddressExists(email string) (bool, error)
	CreateUser(newUser auth.User) error
}

type userService struct {
	log   ilog.StdLogger
	query repository.ApiUserQuery
}

func NewUserService(apiUserQuery repository.ApiUserQuery, logger ilog.StdLogger) UserService {
	return &userService{
		log:   logger,
		query: apiUserQuery,
	}
}

func (s *userService) GetUserByEmailAddress(email string) (user auth.User, err error) {
	user, err = s.query.GetUserByEmailAddress(email)
	return
}

func (s *userService) UserWithEmailAddressExists(email string) (bool, error) {
	return s.query.UserWithEmailAddressExists(email)
}

func (s *userService) CreateUser(user auth.User) error {
	return s.query.CreateUser(user)
}
