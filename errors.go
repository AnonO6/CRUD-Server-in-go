package main

import "errors"

var (
	ErrUserExists = errors.New("user with this email already exists")
)