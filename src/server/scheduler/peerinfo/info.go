package peerinfo

import (
	"fmt"
	"time"
)

func New(peer string) *PeerInfo {
	return &PeerInfo{
		Name:                 peer,
		Registered:           map[string]bool{},
		OptimisticRegistered: map[string]time.Time{},
	}
}

type PeerInfo struct {
	Name                 string               // Name of peer.
	State                bool                 // Live or dead.
	OptimisticState      bool                 // Live or dead.
	Registered           map[string]bool      // Actors on this peer.
	OptimisticRegistered map[string]time.Time // Possible actors on this peer.
}

func (pi *PeerInfo) NumActors() int {
	return len(pi.Registered) + len(pi.OptimisticRegistered)
}

func (pi *PeerInfo) String() string {
	return fmt.Sprintf(
		`Name:%v State:%v OptimisticState:%v Registered:[%+v] OptimisticRegistered:[%+v]`,
		pi.Name, pi.State, pi.OptimisticState, pi.Registered, pi.OptimisticRegistered,
	)
}
