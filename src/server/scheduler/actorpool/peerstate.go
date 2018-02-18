package actorpool

import (
	"math"
	"sync"
	"time"

	"github.com/epsniff/gui/src/server/scheduler/peerinfo"
)

func newPeersState() *PeersState {
	return &PeersState{
		mu:                   &sync.RWMutex{},
		peers:                map[string]*peerinfo.PeerInfo{},
		registered:           map[string]string{},
		optimisticRegistered: map[string]string{},
	}
}

// PeersState contains state of all the known and optimisiticaly started actors
// and which peers they have been started on.
type PeersState struct {
	mu *sync.RWMutex

	peername string
	// peername --> PeerInfo
	peers map[string]*peerinfo.PeerInfo

	//inverted maps (actorname --> peername)
	registered           map[string]string
	optimisticRegistered map[string]string
}

// State retrieves peer state.
func (ps *PeersState) State(peer string) (state bool, optimisticstate bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return false, false
	}
	state, optimisticstate = pi.State, pi.OptimisticState
	return state, optimisticstate
}

// live peer.
func (ps *PeersState) live(peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.State = live
	pi.OptimisticState = live
}

// optimisticallyLive until an event marks the peer dead.
// Currently this has no affect on scheduling.
func (ps *PeersState) optimisticallyLive(peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.OptimisticState = live
}

// dead peer.
func (ps *PeersState) dead(peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.State = dead
	pi.OptimisticState = dead
}

// optimisticallyDead until a real event marks the peer alive again.
// Making the peer optimistically dead will remove it from any
// scheduling, in other words, it will never be returned as a
// peer from MinAssigned when marked optimistically dead.
func (ps *PeersState) optimisticallyDead(peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.OptimisticState = dead
}

func (ps *PeersState) register(actor, peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		// TODO is this an invariance ?  should we allow a registration for an unknow peer?
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if registered is called twice without a
	// unregister inbetween.
	if oldpeer, ok := ps.registered[actor]; ok {
		oldPi, ok := ps.peers[oldpeer]
		if !ok {
			panic("optimisticRegistered pointed to a non existing peer name???")
		}
		delete(oldPi.Registered, actor)
		delete(ps.registered, actor)
	}
	// Check if the actor was optimistically assigned
	// to a peer. Since this method represents a REAL
	// registration, it overrides any optimistic
	// registration.
	if oldpeer, ok := ps.optimisticRegistered[actor]; ok {
		oldPi, ok := ps.peers[oldpeer]
		if !ok {
			panic("optimisticRegistered pointed to a non existing peer name???")
		}
		delete(oldPi.OptimisticRegistered, actor)
		delete(ps.optimisticRegistered, actor)
	}

	// Update the set of actors for
	// this peer.
	pi.Registered[actor] = true

	// Update the global actor to peer
	// mapping.
	ps.registered[actor] = pi.Name
}

// unregister the actor from its current peer.
func (ps *PeersState) unregister(actor string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !isValidName(actor) {
		return
	}

	// Remove optimistic registrations, since
	// this is a REAL unregister, the actor
	// is for sure not running.
	if opeer, ok := ps.optimisticRegistered[actor]; ok {
		oPi, ok := ps.peers[opeer]
		if !ok {
			panic("optimisticRegistered pointed to a non existing peer name???")
		}
		delete(oPi.OptimisticRegistered, actor)
		delete(ps.optimisticRegistered, actor)
	}

	peer, ok := ps.registered[actor]
	if !ok {
		return
	}

	pi, ok := ps.peers[peer]
	if !ok {
		panic("optimisticRegistered pointed to a non existing peer name???")
	}

	// Remove registrations.
	if _, ok = ps.registered[actor]; !ok {
		// Never registered.
		return
	}
	delete(pi.Registered, actor)
	delete(ps.registered, actor)
}

func (ps *PeersState) optimisticallyRegister(actor, peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}

	// CRITICAL
	// Check if this actor is already registered
	// under a different peer. This could happen
	// if optimistic-registered is called twice
	// without a unregister inbetween.
	if oldpeer, ok := ps.optimisticRegistered[actor]; ok {
		oldPi, ok := ps.peers[oldpeer]
		if !ok {
			panic("optimisticRegistered pointed to a non existing peer name???")
		}
		delete(oldPi.OptimisticRegistered, actor)
		delete(ps.optimisticRegistered, actor)
	}

	// Update the set of actors for
	// this peer.
	pi.OptimisticRegistered[actor] = time.Now()

	// Update the global actor to peer
	// mapping.
	ps.optimisticRegistered[actor] = pi.Name
}

// optimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ps *PeersState) optimisticallyUnregister(actor string) {
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
func (ps *PeersState) Registered() map[string]string {
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
func (ps *PeersState) IsRegistered(actorName string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.registered[actorName]
	return ok
}

// OptimisticallyRegister returns a copy of the peerstate's map of actorname --> peername.
func (ps *PeersState) OptimisticallyRegister() map[string]string {
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
func (ps *PeersState) IsOptimisticallyRegister(actor string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.optimisticRegistered[actor]
	return ok
}

// NumRegistered returns the current number of actors
// registered.
func (ps *PeersState) NumRegistered() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.registered)
}

// NumOptimisticallyRegistered returns the current number
// of actors optimistically registered.
func (ps *PeersState) NumOptimisticallyRegistered() int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return len(ps.optimisticRegistered)
}

// NumRegisteredOn the peer.
func (ps *PeersState) NumRegisteredOn(peer string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.Registered)
}

// NumOptimisticallyRegisteredOn the peer.
func (ps *PeersState) NumOptimisticallyRegisteredOn(peer string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.OptimisticRegistered)
}

//returns a struct that represents this peer queue's internal state.  Used for loggging.
func (ps *PeersState) Status() *PeersStatus {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	//TODO add states like peer.burden, peer.actorCnt, pq.NumbActor
	alive := len(ps.registered)
	if alive <= 0 {
		return nil
	}

	pcnt := 0
	for _, p := range ps.peers {
		if p.State == live {
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
