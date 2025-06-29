package errors

import "errors"

var ErrEmailAlreadyExist = errors.New("email already exist")

var ErrEmailOrPassWordFalse = errors.New("email/password is wrong")

var ErrIdNotFound = errors.New("Id not found")

var ErrPackageNotFound = errors.New("Package not found")
