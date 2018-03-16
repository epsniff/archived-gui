package info

import (
	"fmt"
)

func New(peer string) *PeerInfo {
	return &PeerInfo{
		Name: peer,
	}
}

type PeerInfo struct {
	Name            string // Name of peer.
	State           bool   // Live or dead.
	OptimisticState bool   // Live or dead.
}

func (pi *PeerInfo) String() string {
	return fmt.Sprintf(
		`Name:%v State:%v OptimisticState:%v`,
		pi.Name, pi.State, pi.OptimisticState,
	)
}
