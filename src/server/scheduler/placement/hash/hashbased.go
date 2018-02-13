package hash

import (
	"hash/fnv"
	"math"
)

// Relocate actors from peers that have an unfair number of
// actors, where unfair is defined as the integer ceiling
// of the average.
func (pq *ActorPool) Relocate() *RelocationPlan {
	if !pq.rebalance {
		return nil
	}

	pq.mu.Lock()
	defer pq.mu.Unlock()

	alive := len(pq.registered)
	if alive <= 0 {
		return nil
	}

	pcnt := 0
	for _, p := range pq.peers {
		if p.state == live {
			pcnt++
		}
	}
	if pcnt <= 0 {
		return nil
	}

	avePerPeer := int(math.Ceil(float64(alive) / float64(pcnt)))
	if avePerPeer < 0 {
		avePerPeer = 1
	}

	plan := NewRelocationPlan(pq.actorType, alive, avePerPeer)
	for name, p := range pq.peers {
		// REVIEWER NOTES: with this block commmented out we'll rebalance all actors on a dead node.
		//    and the calls to pq.byHash will already deal with the dead peers.
		// if p.state == dead {
		// 	plan.Burden[name] = 0
		// 	continue
		// }
		burden := len(p.registered) - avePerPeer
		plan.Peers = append(plan.Peers, name)
		plan.Count[name] = p.NumActors()
		plan.Burden[name] = burden

		for actor := range p.registered {
			if peer, err := pq.byHash(actor); err != nil {
				//TODO log error
				continue
			} else if peer != name {
				plan.Relocations = append(plan.Relocations, actor)
			}
		}
	}

	return plan
}

// ByHash based selection of "next" living peer
// based on name.
func (pq *ActorPool) BestPeer(name string) (string, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.byHash(name)
}

func (pq *ActorPool) byHash(name string) (string, error) {
	if len(pq.activePeers.peers) == 0 {
		return "", ErrEmpty
	}

	//TODO return ErrEmpty if len(pq.selector.peer) == 0???

	h := fnv.New64()
	h.Write([]byte(name))
	v := h.Sum64()
	l := uint64(len(pq.activePeers.peers))
	peer := pq.activePeers.peers[int(v%l)]
	return peer, nil
}
