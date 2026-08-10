package engine

import "errors"

var (
	ErrCurrentNodeNotFound = errors.New("current node not found")
	ErrChoiceNotFound      = errors.New("choice not found in current node")
	ErrInvalidChoiceTarget = errors.New("invalid choice target")
)
