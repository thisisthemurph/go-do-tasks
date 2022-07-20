package api

import "errors"

var (
	ErrorAccountNotFound      = errors.New("the specified account could not be found")
	ErrorAccountNotCreated    = errors.New("the account could not be created")
	ErrorAccountAlreadyExists = errors.New("the specified account already exists")
)

var (
	ErrorUserNotFound       = errors.New("a user with the specified email address could not be found")
	ErrorUserAlreadyExists  = errors.New("a username with the given email address already exists")
	ErrorUserAuthentication = errors.New("a user with the given email and password combination could not be found")
)

var (
	ErrorProjectNotFound   = errors.New("the requested project could not be found")
	ErrorProjectNotCreated = errors.New("the project could not be created")
	ErrorProjectJsonParse  = errors.New("could not process the given project")
)

var (
	ErrorStoryNotFound   = errors.New("the specified story could not be found")
	ErrorStoryNotCreated = errors.New("the story could not be created")
)
