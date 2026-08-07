package attempt

import "errors"

var (
	ErrAttemptNotFound         = errors.New("attempt not found")
	ErrActiveAttemptNotFound   = errors.New("active attempt not found")
	ErrActiveAttemptExists     = errors.New("active attempt already exists")
	ErrAttemptNotInProgress    = errors.New("attempt is not in progress")
	ErrInvalidStatusTransition = errors.New("invalid attempt status transition")
	ErrInvalidAttemptState     = errors.New("invalid attempt state")
)