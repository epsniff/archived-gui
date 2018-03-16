package pool

import (
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
	"github.com/epsniff/gui/src/server/scheduler/types"
	"github.com/lytics/grid"
)

// New
func New(ps clusterstate.PeersState) *ActorPool {
	return &ActorPool{
		actorPoolState: newActorPoolState(ps),
	}
}

type ActorPool struct {
	actorPoolState *ActorPoolState
}

func (ap *ActorPool) Remove(name string) {
	ap.actorPoolState.remove(name)
}

func (ap *ActorPool) IsRequired(name string) bool {
	return ap.actorPoolState.isRequired(name)
}

// Missing actors that are required but not registered.
func (ap *ActorPool) Missing() []*grid.ActorStart {
	return ap.actorPoolState.missing()
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumRegistered() int {
	return ap.actorPoolState.NumRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumRegisteredOn(peer string) int {
	return ap.actorPoolState.NumRegisteredOn(peer)
}

//NumRegistered returns the number of actors registered in this pool
func (ap *ActorPool) NumOptimisticallyRegistered() int {
	return ap.actorPoolState.NumOptimisticallyRegistered()
}

//NumRegisteredOn returns the number of actors registered in a peer in the pool
func (ap *ActorPool) NumOptimisticallyRegisteredOn(peer string) int {
	return ap.actorPoolState.NumOptimisticallyRegisteredOn(peer)
}

//IsRegistered returns if the actorName as been registered already.
func (ap *ActorPool) IsRegistered(actorName string) bool {
	return ap.actorPoolState.IsRegistered(actorName)
}

func (ap *ActorPool) IsOptimisticallyRegistered(actorName string) bool {
	return ap.actorPoolState.IsOptimisticallyRegister(actorName)
}

// Actors that are registered in this pool
func (ap *ActorPool) Actors() []*ActorEntry {
	return ap.actorPoolState.getActors()
}

// Register the actor to the peer.
func (ap *ActorPool) Register(isRequired bool, actor *grid.ActorStart, peer string) error {
	if !isValidName(actor.Name) {
		return types.ErrInvalidActorName
	}
	if !isValidName(peer) {
		return types.ErrInvalidPeerName
	}
	return ap.actorPoolState.register(isRequired, actor, peer)
}

// OptimisticallyRegister an actor, ie: no confirmation has
// arrived that the actor is actually running on the peer,
// but it has been requested to run on the peer.
func (ap *ActorPool) OptimisticallyRegister(isRequired bool, actor *grid.ActorStart, peer string) error {
	if !isValidName(actor.Name) {
		return types.ErrInvalidActorName
	}
	if !isValidName(peer) {
		return types.ErrInvalidPeerName
	}
	return ap.actorPoolState.optimisticallyRegister(isRequired, actor, peer)
}

// Unregister the actor from its current peer.
func (ap *ActorPool) Unregister(actor string) error {
	if !isValidName(actor) {
		return types.ErrInvalidActorName
	}
	return ap.actorPoolState.unregister(actor)
}

// OptimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ap *ActorPool) OptimisticallyUnregister(actor string) error {
	if !isValidName(actor) {
		return types.ErrInvalidActorName
	}
	return ap.actorPoolState.optimisticallyUnregister(actor)
}

/*
//Status returns a struct that represents all the peer queue's internal states used
//for logging and debugging
func (ap *ActorPool) Status() *PeersStatus {
	return ap.actorPoolState.Status()
}
*/
func isValidName(name string) bool {
	return name != ""
}
