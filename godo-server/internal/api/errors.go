package api

import "errors"

var ProjectNotFoundError = errors.New("The requested project could not be found")
var ProjectNotCreatedError = errors.New("The project could not be created")
