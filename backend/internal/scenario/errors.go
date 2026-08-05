package scenario

import "errors"

var (
	ErrEmptyStartNodeID       = errors.New("empty start node ID")
	ErrUnknownStartNode       = errors.New("unknown start node")
	ErrEmptyNodeID            = errors.New("empty node ID")
	ErrDuplicateNodeID        = errors.New("duplicate node ID")
	ErrEmptyEndingID          = errors.New("empty ending ID")
	ErrDuplicateEndingID      = errors.New("duplicate ending ID")
	ErrEmptyChoiceID          = errors.New("empty choice ID")
	ErrDuplicateChoiceID      = errors.New("duplicate choice ID")
	ErrInvalidChoiceTarget    = errors.New("choice must have exactly one target")
	ErrUnknownNode            = errors.New("unknown node")
	ErrUnknownEnding          = errors.New("unknown ending")
	ErrInvalidWeight          = errors.New("invalid weight")
	ErrInvalidScore           = errors.New("invalid score")
	ErrMissingRiskCategory    = errors.New("missing risk category")
	ErrCycleDetected          = errors.New("cycle detected")
	ErrUnreachableNode        = errors.New("unreachable node")
	ErrUnreachableEnding      = errors.New("unreachable ending")
	ErrTooFewReachableEndings = errors.New("too few reachable endings")
)
