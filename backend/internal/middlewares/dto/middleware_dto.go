package dto

import "errors"

const (
	FAILED_AUTH = "failed_auth"
)

var (
	ErrHeaderMissing = errors.New("Authorization header missing")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidHeader = errors.New("invalid authentication header")
)
