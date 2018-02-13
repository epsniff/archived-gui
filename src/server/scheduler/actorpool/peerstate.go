package actorpool

import (
	"sync"
	"time"

	"github.com/epsniff/gui/src/server/scheduler/peerinfo"
)

func newPeerState(name string) *PeersState {
	return &PeersState{
		mu:                       &sync.RWMutex{},
		peers:                    map[string]*peerinfo.PeerInfo{},
		registered:               map[string]string{},
		optimisticallyRegistered: map[string]string{},
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
	registered               map[string]string
	optimisticallyRegistered map[string]string
}

func (ps *PeersState) register(actor, peer string) {
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
	// if registered is called twice without a
	// unregister inbetween.
	if pi, ok := ps.registered[actor]; ok {
		delete(pi.Registered, actor)
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
}

func (ps *PeersState) optimisticallyRegister(actor, peer string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

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
}

// optimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ps *PeersState) optimisticallyUnregister(actor string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := pq.optimisticRegistered[actor]
	if !ok {
		// Never registered.
		return
	}
	delete(pi.optimisticRegistered, actor)
	delete(pq.optimisticRegistered, actor)
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
	for k, v := range ps.optimisticallyRegistered {
		res[k] = v
	}
	return res
}

// IsOptimisticallyRegister returns true if the actor
// has been optimistically registered.
func (ps *PeersState) IsOptimisticallyRegister(actorame string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	_, ok := ps.optimisticallyRegistered[actorName]
	return ok
}

// NumRegistered returns the current number of actors
// registered.
func (ps *PeersState) NumRegistered() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.registered)
}

// NumOptimisticallyRegistered returns the current number
// of actors optimistically registered.
func (ps *PeersState) NumOptimisticallyRegistered() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	return len(ps.optimisticRegistered)
}

// NumRegisteredOn the peer.
func (ps *PeersState) NumRegisteredOn(peer string) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.registered)
}

// NumOptimisticallyRegisteredOn the peer.
func (ps *PeersState) NumOptimisticallyRegisteredOn(peer string) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return 0
	}
	return len(pi.optimisticRegistered)
}
