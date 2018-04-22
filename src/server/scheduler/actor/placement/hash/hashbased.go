package hash

import (
	"hash/fnv"
	"sort"

	"github.com/epsniff/gui/src/server/scheduler/actor/placement/plan"
	"github.com/epsniff/gui/src/server/scheduler/actor/pool"
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
	"github.com/epsniff/gui/src/server/scheduler/types"
)

func New() *HashPlacement {
	return &HashPlacement{}
}

type HashPlacement struct {
}

// Relocate actors from peers that have an unfair number of
// actors, where unfair is defined as the integer ceiling
// of the average.
func (h *HashPlacement) Relocate(peersState clusterstate.PeersState, pools map[string]*pool.ActorPool) (*plan.RelocationPlan, error) {

	relocplan := plan.New()
	if len(pools) == 0 {
		return nil, nil
	}
	peersInfos := peersState.Peers()
	peers := []string{}
	for _, p := range peersInfos {
		peers = append(peers, p.Name)
	}

	sort.Strings(peers)

	for pn, p := range pools {
		for _, actorEntry := range p.Actors() {
			owningPeer, err := h.byHash(actorEntry.Def.Name, peers)
			if err != nil {
				return nil, err
			}
			if actorEntry.Peer != owningPeer {
				relocplan.Relocations = append(
					relocplan.Relocations,
					&plan.Relocation{PoolName: pn, ActorName: actorEntry.Def.Name, Def: actorEntry.Def},
				)
			}
		}
	}

	return relocplan, nil
}

// ByHash based selection of "next" living peer
// based on name.
func (h *HashPlacement) BestPeer(actorName string, peersState clusterstate.PeersState, pool map[string]*pool.ActorPool) (string, error) {
	peersInfos := peersState.Peers()
	peers := []string{}
	for _, p := range peersInfos {
		peers = append(peers, p.Name)
	}

	sort.Strings(peers)

	return h.byHash(actorName, peers)
}

func (h *HashPlacement) byHash(name string, peers []string) (string, error) {
	if len(peers) == 0 {
		return "", types.ErrEmpty
	}

	ha := fnv.New64()
	ha.Write([]byte(name))
	v := ha.Sum64()
	l := uint64(len(peers))
	peer := peers[int(v%l)]
	return peer, nil
}
