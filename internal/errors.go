package internal

import "errors"

var (
	ErrNoKey = errors.New("invalid key:key not found")
)