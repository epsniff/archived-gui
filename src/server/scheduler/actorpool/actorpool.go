package actorpool

import (
	"errors"
	"hash/fnv"
	"math"
	"sync"
	"time"

	"github.com/epsniff/gui/src/server/scheduler/peerinfo"
	"github.com/lytics/grid"
)

const (
	live = true
	dead = false
)

var (
	ErrEmpty             = errors.New("empty")
	ErrInvalidName       = errors.New("invalid name")
	ErrActorTypeMismatch = errors.New("actor type mismatch")
)

// New peer queue.
func New(rebalance bool) *ActorPool {
	return &ActorPool{
		required:             map[string]*grid.ActorStart{},
		registered:           map[string]*peerinfo.PeerInfo{},
		optimisticRegistered: map[string]*peerinfo.PeerInfo{},
		rebalance:            rebalance,
	}
}

type ActorPool struct {
	mu                   sync.Mutex
	actorType            string
	required             map[string]*grid.ActorStart
	registered           map[string]*peerinfo.PeerInfo
	optimisticRegistered map[string]*peerinfo.PeerInfo
	rebalance            bool
}

// IsRequired actor.
func (pq *ActorPool) IsRequired(actor string) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	_, ok := pq.required[actor]
	return ok
}

// ActorType of this peer queue.
func (pq *ActorPool) ActorType() string {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return pq.actorType
}

// SetRequired flag on actor. If it's type does not match
// the type of previously set actors an error is returned.
func (pq *ActorPool) SetRequired(def *grid.ActorStart) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(def.Name) {
		return ErrInvalidName
	}

	if pq.actorType == "" {
		pq.actorType = def.Type
	}
	if pq.actorType != def.Type {
		return ErrActorTypeMismatch
	}
	pq.required[def.Name] = def
	return nil
}

// UnsetRequired flag on actor.
func (pq *ActorPool) UnsetRequired(actor string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	delete(pq.required, actor)
	if len(pq.required) == 0 {
		pq.actorType = ""
	}
}

// Missing actors that are required but not registered.
func (pq *ActorPool) Missing() []*grid.ActorStart {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	var missing []*grid.ActorStart
	for name, def := range pq.required {
		if _, ok := pq.registered[name]; !ok {
			missing = append(missing, def)
		}
	}
	return missing
}

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
func (pq *ActorPool) ByHash(name string) (string, error) {
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

// IsRegistered returns true if the actor has been
// registered.
func (pq *ActorPool) IsRegistered(actor string) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	_, ok := pq.registered[actor]
	return ok
}

// IsOptimisticallyRegistered returns true if the actor
// has been optimistically registered.
func (pq *ActorPool) IsOptimisticallyRegistered(actor string) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	_, ok := pq.optimisticRegistered[actor]
	return ok
}

// NumRegistered returns the current number of actors
// registered.
func (pq *ActorPool) NumRegistered() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return len(pq.registered)
}

// NumOptimisticallyRegistered returns the current number
// of actors optimistically registered.
func (pq *ActorPool) NumOptimisticallyRegistered() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return len(pq.optimisticRegistered)
}

// NumRegisteredOn the peer.
func (pq *ActorPool) NumRegisteredOn(peer string) int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pi, ok := pq.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.registered)
}

// NumOptimisticallyRegisteredOn the peer.
func (pq *ActorPool) NumOptimisticallyRegisteredOn(peer string) int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pi, ok := pq.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.optimisticRegistered)
}

// Register the actor to the peer.
func (pq *ActorPool) Register(actor, peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(actor) {
		return
	}
	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if registered is called twice without a
	// unregister inbetween.
	if pi, ok := pq.registered[actor]; ok {
		delete(pi.registered, actor)
		delete(pq.registered, actor)
	}
	// Check if the actor was optimistically assigned
	// to a peer. Since this method represents a REAL
	// registration, it overrides any optimistic
	// registration.
	if pi, ok := pq.optimisticRegistered[actor]; ok {
		delete(pi.optimisticRegistered, actor)
		delete(pq.optimisticRegistered, actor)
	}

	// Update the set of actors for
	// this peer.
	pi.registered[actor] = true

	// Update the global actor to peer
	// mapping.
	pq.registered[actor] = pi

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

// OptimisticallyRegister an actor, ie: no confirmation has
// arrived that the actor is actually running on the peer,
// but it has been requested to run on the peer.
func (pq *ActorPool) OptimisticallyRegister(actor, peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(actor) {
		return
	}
	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if optimistic-registered is called twice
	// without a unregister inbetween.
	if pi, ok := pq.optimisticRegistered[actor]; ok {
		delete(pi.optimisticRegistered, actor)
		delete(pq.optimisticRegistered, actor)
	}

	// Update the set of actors for
	// this peer.
	pi.optimisticRegistered[actor] = time.Now()

	// Update the global actor to peer
	// mapping.
	pq.optimisticRegistered[actor] = pi

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

// Unregister the actor from its current peer.
func (pq *ActorPool) Unregister(actor string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(actor) {
		return
	}

	// Remove optimistic registrations, since
	// this is a REAL unregister, the actor
	// is for sure not running.
	pi, ok := pq.optimisticRegistered[actor]
	if ok {
		delete(pi.optimisticRegistered, actor)
		delete(pq.optimisticRegistered, actor)
	}

	// Remove registrations.
	pi, ok = pq.registered[actor]
	if !ok {
		// Never registered.
		return
	}
	delete(pi.registered, actor)
	delete(pq.registered, actor)

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

// OptimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (pq *ActorPool) OptimisticallyUnregister(actor string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(actor) {
		return
	}

	pi, ok := pq.optimisticRegistered[actor]
	if !ok {
		// Never registered.
		return
	}
	delete(pi.optimisticRegistered, actor)
	delete(pq.optimisticRegistered, actor)

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

func isValidName(name string) bool {
	return name != ""
}
