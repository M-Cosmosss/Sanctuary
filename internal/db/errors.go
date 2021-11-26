package db

import "github.com/pkg/errors"

var ErrRouteAlreadyExists = errors.New("Route Path already exists.")
var ErrRouteNotExists = errors.New("Route Path does not exist.")

var ErrRouteGroupAlreadyExists = errors.New("Route Group Name or Path already exists.")
var ErrRouteGroupNotExists = errors.New("Route Group does not exist.")

var ErrServiceAlreadyExists = errors.New("Service already exists.")
var ErrServiceNotExists = errors.New("Service does not exist.")

var ErrServiceGroupAlreadyExists = errors.New("Service Group Name already exists.")
var ErrServiceGroupNotExists = errors.New("Service Group does not exist.")

var ErrUnknown = errors.New("Unknown database error.")

var ErrNotHTTPMethod = errors.New("Method is not in the HTTP")
