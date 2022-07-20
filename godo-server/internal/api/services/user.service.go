package services

import (
	"godo/internal/api"
	"godo/internal/helper/ilog"
	"godo/internal/repository"
	"godo/internal/repository/entities"
)

type UserService interface {
	GetUserByEmailAddress(email string) (user *entities.User, err error)
	UserWithEmailAddressExists(email string) (bool, error)
	CreateUser(newUser entities.User) (*entities.User, error)
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

func (s *userService) GetUserByEmailAddress(email string) (*entities.User, error) {
	exists, err := s.query.UserWithEmailAddressExists(email)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, api.ErrorUserNotFound
	}

	var user *entities.User
	user, err = s.query.GetUserByEmailAddress(email)
	return user, err
}

func (s *userService) CreateUser(newUser entities.User) (*entities.User, error) {
	exists, err := s.query.UserWithEmailAddressExists(newUser.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, api.ErrorUserAlreadyExists
	}

	return s.query.CreateUser(newUser)
}

func (s *userService) UserWithEmailAddressExists(email string) (bool, error) {
	return s.query.UserWithEmailAddressExists(email)
}
