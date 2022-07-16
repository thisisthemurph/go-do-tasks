package api

import "errors"

var ProjectNotFoundError = errors.New("The requested project could not be found")
var ProjectNotCreatedError = errors.New("The project could not be created")

var UserNotFoundError = errors.New("A user with the specified email address could not be found")
var UserAlreadyExistsError = errors.New("A username with the given email address already exists")
