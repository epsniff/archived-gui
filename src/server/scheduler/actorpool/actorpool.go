package actorpool

import (
	"sync"

	"github.com/lytics/grid"
)

// New peer queue.
func New(rebalance bool) *ActorPool {
	return &ActorPool{
		required:  map[string]*grid.ActorStart{},
		peerState: newPeersState(),
		rebalance: rebalance,
	}
}

type ActorPool struct {
	mu        sync.RWMutex
	required  map[string]*grid.ActorStart
	actorType string
	rebalance bool

	peerState *PeersState
}

// IsRequired actor.
func (ap *ActorPool) IsRequired(actor string) bool {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	_, ok := ap.required[actor]
	return ok
}

// ActorType of this peer queue.
func (ap *ActorPool) ActorType() string {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	return ap.actorType
}

// SetRequired flag on actor. If it's type does not match
// the type of previously set actors an error is returned.
func (ap *ActorPool) SetRequired(def *grid.ActorStart) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if !isValidName(def.Name) {
		return ErrInvalidName
	}

	if ap.actorType == "" {
		ap.actorType = def.Type
	}
	if ap.actorType != def.Type {
		return ErrActorTypeMismatch
	}
	ap.required[def.Name] = def
	return nil
}

// UnsetRequired flag on actor.
func (ap *ActorPool) UnsetRequired(actor string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	delete(ap.required, actor)
	if len(ap.required) == 0 {
		ap.actorType = ""
	}
}

// Missing actors that are required but not registered.
func (ap *ActorPool) Missing() []*grid.ActorStart {
	var missing []*grid.ActorStart
	for name, def := range ap.required {
		if isReg := ap.peerState.IsRegistered(name); !isReg {
			missing = append(missing, def)
		}
	}
	return missing
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumRegistered() int {
	return ap.peerState.NumRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumRegisteredOn(peer string) int {
	return ap.peerState.NumRegisteredOn(peer)
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumOptimisticallyRegistered() int {
	return ap.peerState.NumOptimisticallyRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumOptimisticallyRegisteredOn(peer string) int {
	return ap.peerState.NumOptimisticallyRegisteredOn(peer)
}

//IsRegistered returns if the actorName as been registered already.
func (ap *ActorPool) IsRegistered(actorName string) bool {
	return ap.peerState.IsRegistered(actorName)
}

func (ap *ActorPool) IsOptimisticallyRegistered(actorName string) bool {
	return ap.peerState.IsOptimisticallyRegister(actorName)
}

// Register the actor to the peer.
func (ap *ActorPool) Register(actor, peer string) {
	if !isValidName(actor) {
		return
	}
	if !isValidName(peer) {
		return
	}
	ap.peerState.register(actor, peer)
}

// OptimisticallyRegister an actor, ie: no confirmation has
// arrived that the actor is actually running on the peer,
// but it has been requested to run on the peer.
func (ap *ActorPool) OptimisticallyRegister(actor, peer string) {
	if !isValidName(actor) {
		return
	}
	if !isValidName(peer) {
		return
	}
	ap.peerState.optimisticallyRegister(actor, peer)
}

// Unregister the actor from its current peer.
func (ap *ActorPool) Unregister(actor string) {
	if !isValidName(actor) {
		return
	}
	ap.peerState.unregister(actor)
}

// OptimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ap *ActorPool) OptimisticallyUnregister(actor string) {
	if !isValidName(actor) {
		return
	}
	ap.peerState.optimisticallyUnregister(actor)
}

// Live peer.
func (ap *ActorPool) Live(peer string) {
	if !isValidName(peer) {
		return
	}
	ap.peerState.live(peer)
}

// OptimisticallyLive until an event marks the peer dead.
// Currently this has no affect on scheduling.
func (ap *ActorPool) OptimisticallyLive(peer string) {
	if !isValidName(peer) {
		return
	}
	ap.peerState.optimisticallyLive(peer)
}

// Dead peer.
func (ap *ActorPool) Dead(peer string) {
	if !isValidName(peer) {
		return
	}
	ap.peerState.dead(peer)
}

// OptimisticallyDead until a real event marks the peer alive again.
// Making the peer optimistically dead will remove it from any
// scheduling, in other words, it will never be returned as a
// peer from MinAssigned when marked optimistically dead.
func (ap *ActorPool) OptimisticallyDead(peer string) {
	if !isValidName(peer) {
		return
	}
	ap.peerState.optimisticallyDead(peer)
}

//Status returns a struct that represents all the peer queue's internal states used
//for logging and debugging
func (ap *ActorPool) Status() *PeersStatus {
	return ap.peerState.Status()
}

func isValidName(name string) bool {
	return name != ""
}
