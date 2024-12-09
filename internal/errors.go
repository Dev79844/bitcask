package internal

import "errors"

var (
	ErrNoKey 	=	errors.New("invalid key:key not found")
	ErrEmptyKey	= errors.New("invalid key: key cannot be empty")
)