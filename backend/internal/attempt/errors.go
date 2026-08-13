package attempt

import "errors"

var (
	ErrAttemptNotFound          = errors.New("attempt not found")
	ErrActiveAttemptNotFound    = errors.New("active attempt not found")
	ErrCompletedAttemptNotFound = errors.New("completed attempt not found")
	ErrActiveAttemptExists      = errors.New("active attempt already exists")
	ErrAttemptNotInProgress     = errors.New("attempt is not in progress")
	ErrInvalidStatusTransition  = errors.New("invalid attempt status transition")
	ErrInvalidAttemptState      = errors.New("invalid attempt state")
	ErrIdempotencyConflict      = errors.New("idempotency key conflicts with request payload")
	ErrNodeAlreadyAnswered      = errors.New("attempt node already answered")
	ErrAttemptNodeMismatch      = errors.New("attempt current node does not match request node")
	ErrAnswerNotFound           = errors.New("answer not found")
)
