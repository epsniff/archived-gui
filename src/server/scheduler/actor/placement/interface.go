package placement

import (
	"github.com/epsniff/gui/src/server/scheduler/actor/placement/plan"
	"github.com/epsniff/gui/src/server/scheduler/actor/pool"
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
)

type Placement interface {
	Relocate(peers clusterstate.PeersState, pools map[string]*pool.ActorPool) (*plan.RelocationPlan, error)
	BestPeer(actorName string, peers clusterstate.PeersState, pools map[string]*pool.ActorPool) (string, error)
}
