package api

import "errors"

var (
	AccountNotFoundError      = errors.New("The specified account could not be found")
	AccountNotCreatedError    = errors.New("The account could not be created")
	AccountAlreadyExistsError = errors.New("The specified account already exists")
)

var (
	UserNotFoundError       = errors.New("A user with the specified email address could not be found")
	UserAlreadyExistsError  = errors.New("A username with the given email address already exists")
	UserAuthenticationError = errors.New("A user with the given email and password combination could not be found")
)

var (
	ProjectNotFoundError   = errors.New("The requested project could not be found")
	ProjectNotCreatedError = errors.New("The project could not be created")
	ProjectJsonParseError  = errors.New("Could not process the given project")
)
