package placement

import "github.com/epsniff/gui/src/server/scheduler/actorpool"

type Placement interface {
	Relocate(peers actorpool.PeersState) *RelocationPlan
	BestPeer(actorName string, peers actorpool.PeersState) (string, error)
}
