package actorpool

import "errors"

const (
	live = true
	dead = false
)

var (
	ErrEmpty             = errors.New("empty")
	ErrInvalidName       = errors.New("invalid name")
	ErrActorTypeMismatch = errors.New("actor type mismatch")
)
