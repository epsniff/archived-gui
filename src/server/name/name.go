package name

import (
	"errors"
)

const (
	LeaderActor            = "leader"
	PeerMonitorActorPrefix = "peer-monitor"
)

var (
	ErrUnknownActorType = errors.New("unknown actor type")
)
