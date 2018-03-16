package pool

import (
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
	"github.com/epsniff/gui/src/server/scheduler/types"
	"github.com/lytics/grid"
)

func newActorPoolState(ps clusterstate.PeersState) *ActorPoolState {
	return &ActorPoolState{
		mu:                   &sync.RWMutex{},
		peers:                ps,
		registered:           map[string]string{},
		optimisticRegistered: map[string]string{},
		actors:               map[string]*ActorEntry{},
	}
}

type ActorEntry struct {
	Def          *grid.ActorStart
	IsRequired   bool
	IsRegistered bool //registered vs opimiamisically registered
	Peer         string
}

// ActorPoolState contains state of all the known and optimisiticaly started actors
// and which peers they have been started on.
type ActorPoolState struct {
	mu *sync.RWMutex

	peers clusterstate.PeersState

	//inverted maps (actorname --> peername)
	registered           map[string]string
	optimisticRegistered map[string]string

	//actor name --> {isRequired, ActorDef}
	actors map[string]*ActorEntry
}

func (ps *ActorPoolState) missing() []*grid.ActorStart {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	var missing []*grid.ActorStart

	for name, entry := range ps.actors {
		if !entry.IsRequired {
			continue
		}
		if !ps.IsRegistered(name) {
			missing = append(missing, entry.Def)
		}
	}
	return missing
}

func (ps *ActorPoolState) remove(name string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.actors, name)
}

func (ps *ActorPoolState) isRequired(name string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	e, ok := ps.actors[name]
	if !ok {
		return false
	}
	return e.IsRequired
}

func (ps *ActorPoolState) getActors() []*ActorEntry {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	as := []*ActorEntry{}
	for _, ae := range ps.actors {
		as = append(as, ae)
	}
	return as
}

func (ps *ActorPoolState) register(isRequired bool, actor *grid.ActorStart, peer string) error {
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
	delete(ps.registered, actor.Name)

	// Check if the actor was optimistically assigned
	// to a peer. Since this method represents a REAL
	// registration, it overrides any optimistic
	// registration.
	delete(ps.optimisticRegistered, actor.Name)

	// Update the global actor to peer
	// mapping.
	ps.registered[actor.Name] = pi.Name

	ps.actors[actor.Name] = &ActorEntry{actor, isRequired, true, peer}

	return nil
}

// unregister the actor from its current peer.
func (ps *ActorPoolState) unregister(actor string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !isValidName(actor) {
		return types.ErrInvalidActorName
	}

	// Remove optimistic registrations, since
	// this is a REAL unregister, the actor
	// is for sure not running.
	delete(ps.optimisticRegistered, actor)

	_, ok := ps.registered[actor]
	if !ok {
		return types.ErrActorNotRegistered
	}

	// Remove registrations.
	delete(ps.registered, actor)
	return nil
}

func (ps *ActorPoolState) optimisticallyRegister(isRequired bool, actor *grid.ActorStart, peer string) error {
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
	delete(ps.optimisticRegistered, actor.Name)

	// Update the global actor to peer
	// mapping.
	ps.optimisticRegistered[actor.Name] = pi.Name

	ps.actors[actor.Name] = &ActorEntry{actor, isRequired, false, peer}

	return nil
}

// optimisticallyUnregister the actor, ie: no confirmation has
// arrived that the actor is NOT running on the peer, but
// perhaps because of a failed request to the peer to start
// the actor it is known that likely the actor is not running.
func (ps *ActorPoolState) optimisticallyUnregister(actor string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.optimisticRegistered, actor)
	return nil
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

	cnt := 0
	for _, p := range ps.registered {
		if p == peer {
			cnt++
		}
	}
	return cnt
}

// NumOptimisticallyRegisteredOn the peer.
func (ps *ActorPoolState) NumOptimisticallyRegisteredOn(peer string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	cnt := 0
	for _, p := range ps.optimisticRegistered {
		if p == peer {
			cnt++
		}
	}
	return cnt
}

/*
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
	for _, p := range ps.peers.Get {
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
*/
