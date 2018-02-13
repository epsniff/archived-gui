package actorpool

import (
	"sync"

	"github.com/lytics/grid"
)

// New peer queue.
func New(rebalance bool) *ActorPool {
	return &ActorPool{
		required:  map[string]*grid.ActorStart{},
		peerState: newPeerState(),
		rebalance: rebalance,
	}
}

type ActorPool struct {
	mu        sync.Mutex
	actorType string
	rebalance bool

	peerState *PeersState
	required  map[string]*grid.ActorStart
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
		if isReg := pq.peerState.IsRegistered(name); !isReg {
			missing = append(missing, def)
		}
	}
	return missing
}

func (pq *ActorPool) IsRegistered(actorName string) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.peerState.IsRegistered(actorName)
}

func (pq *ActorPool) IsOptimisticallyRegistered(actorName string) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.peerState.IsOptimisticallyRegister(actorName)
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
	pq.peerState.register(actor, peer)
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
	pq.peerState.optimisticallyRegister(actor, peer)
}

// Unregister the actor from its current peer.
func (pq *ActorPool) Unregister(actor string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(actor) {
		return
	}
	pq.peerState.unregister(actor)
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
	pq.peerState.optimisticallyUnregister(actor)
}

func isValidName(name string) bool {
	return name != ""
}
