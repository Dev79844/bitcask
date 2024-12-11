package bitcask

import "errors"

var (
	ErrNoKey 	=	errors.New("invalid key:key not found")
	ErrEmptyKey	= errors.New("invalid key: key cannot be empty")
	ErrReadOnly = errors.New("operation not allowed in read only mode")
)