package actorpool

import (
	"math"
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/peerinfo"
	"github.com/epsniff/gui/src/server/scheduler/peerstate"
)

func newActorPoolState(ps peerstate.PeersState) *ActorPoolState {
	return &ActorPoolState{
		mu:                   &sync.RWMutex{},
		peers:                ps,
		registered:           map[string]string{},
		optimisticRegistered: map[string]string{},
	}
}

// ActorPoolState contains state of all the known and optimisiticaly started actors
// and which peers they have been started on.
type ActorPoolState struct {
	mu *sync.RWMutex

	peers peerstate.PeersState

	//inverted maps (actorname --> peername)
	registered           map[string]string
	optimisticRegistered map[string]string
}

func (ps *ActorPoolState) register(actor, peer string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, err := ps.peers.Get(peer)
	if err != nil {
		return err
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if registered is called twice without a
	// unregister inbetween.
	delete(ps.registered, actor)

	// Check if the actor was optimistically assigned
	// to a peer. Since this method represents a REAL
	// registration, it overrides any optimistic
	// registration.
	if oldpeer, ok := ps.optimisticRegistered[actor]; ok {
		delete(ps.optimisticRegistered, actor)
	}

	// Update the global actor to peer
	// mapping.
	ps.registered[actor] = pi.Name

	return nil
}

// unregister the actor from its current peer.
func (ps *ActorPoolState) unregister(actor string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !isValidName(actor) {
		return ErrInvalidActorName
	}

	// Remove optimistic registrations, since
	// this is a REAL unregister, the actor
	// is for sure not running.
	delete(ps.optimisticRegistered, actor)

	peer, ok := ps.registered[actor]
	if !ok {
		return ErrActorNotRegistered
	}

	pi, err := ps.peers.Get(peer)
	if err != nil {
		return err
	}

	// Remove registrations.
	delete(ps.registered, actor)
	return nil
}

func (ps *ActorPoolState) optimisticallyRegister(actor, peer string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, err := ps.peers.Get(peer)
	if err != nil {
		return err
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if optimistic-registered is called twice
	// without a unregister inbetween.
	delete(ps.optimisticRegistered, actor)

	// Update the global actor to peer
	// mapping.
	ps.optimisticRegistered[actor] = pi.Name

	return nil
}

// optimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ps *ActorPoolState) optimisticallyUnregister(actor string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	peer, ok := ps.optimisticRegistered[actor]
	if !ok {
		// Never registered.
		return
	}

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}

	delete(pi.OptimisticRegistered, actor)
	delete(ps.optimisticRegistered, actor)
}

// Registered returns a copy of the peerstate's map of actorname --> peername.
func (ps *ActorPoolState) Registered() map[string]string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	res := map[string]string{}
	for k, v := range ps.registered {
		res[k] = v
	}
	return res
}

// IsRegistered returns true if the actor has been
// registered.
func (ps *ActorPoolState) IsRegistered(actorName string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.registered[actorName]
	return ok
}

// OptimisticallyRegister returns a copy of the peerstate's map of actorname --> peername.
func (ps *ActorPoolState) OptimisticallyRegister() map[string]string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	res := map[string]string{}
	for k, v := range ps.optimisticRegistered {
		res[k] = v
	}
	return res
}

// IsOptimisticallyRegister returns true if the actor
// has been optimistically registered.
func (ps *ActorPoolState) IsOptimisticallyRegister(actor string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.optimisticRegistered[actor]
	return ok
}

// NumRegistered returns the current number of actors
// registered.
func (ps *ActorPoolState) NumRegistered() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.registered)
}

// NumOptimisticallyRegistered returns the current number
// of actors optimistically registered.
func (ps *ActorPoolState) NumOptimisticallyRegistered() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.optimisticRegistered)
}

// NumRegisteredOn the peer.
func (ps *ActorPoolState) NumRegisteredOn(peer string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.Registered)
}

// NumOptimisticallyRegisteredOn the peer.
func (ps *ActorPoolState) NumOptimisticallyRegisteredOn(peer string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.OptimisticRegistered)
}

//returns a struct that represents this peer queue's internal state.  Used for loggging.
func (ps *ActorPoolState) Status() *PeersStatus {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	//TODO add states like peer.burden, peer.actorCnt, pq.NumbActor
	alive := len(ps.registered)
	if alive <= 0 {
		return nil
	}

	pcnt := 0
	for _, p := range ps.peers {
		if p.State == peerstate.Live {
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

	peers := map[string]*PeerState{}
	for p, pi := range ps.peers {
		ps := &PeerState{
			Name:            pi.Name,
			State:           pi.State,
			OptimisticState: pi.OptimisticState,
		}
		reg := []string{}
		for actor, _ := range pi.Registered {
			reg = append(reg, actor)
		}
		ps.Registered = reg
		optreg := []string{}
		for actor, _ := range pi.Registered {
			optreg = append(optreg, actor)
		}
		ps.OptimisticRegistered = optreg
		ps.Actors = pi.NumActors()
		ps.Burden = len(pi.Registered) - avePerPeer
		peers[p] = ps
	}

	registered := map[string]*ActorState{}
	for p, peer := range ps.registered {
		pi := ps.peers[peer]
		ps := &ActorState{
			Peer:              pi.Name,
			IsAlive:           pi.State,
			IsOptimisticAlive: pi.OptimisticState,
		}
		registered[p] = ps
	}

	optimisticRegistered := map[string]*ActorState{}
	for p, peer := range ps.optimisticRegistered {
		pi := ps.peers[peer]
		ps := &ActorState{
			Peer:              pi.Name,
			IsAlive:           pi.State,
			IsOptimisticAlive: pi.OptimisticState,
		}
		optimisticRegistered[p] = ps
	}
	return &PeersStatus{
		AveActorsPerPeer:     avePerPeer,
		Peers:                peers,
		Registered:           registered,
		OptimisticRegistered: optimisticRegistered,
	}
}
