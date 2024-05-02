package repository

import "errors"

var (
	ErrUserDoesntExist    = errors.New("user does not exist")
	ErrCouldntCreatOrder  = errors.New("couldn't create order")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrWrongProductId     = errors.New("wrong product id ")
)
