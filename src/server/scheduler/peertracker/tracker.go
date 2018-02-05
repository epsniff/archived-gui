package peertracker

import (
	"errors"
	"sort"

	"github.com/epsniff/spider/src/server/scheduler/actorpool"
)

var (
	ErrActorTypeAlreadyRegistered = errors.New("actor type already registered")
	ErrUnknownActorType           = errors.New("unknown actor type")
)

type selector struct {
	activePeers    []string
	peers    map[string]*peerinfo.PeerInfo
}

func New() *Tracker {
	return &Tracker{
		pools: map[string]*actorpool.ActorPool{},
	}
}

type Tracker struct {
	pools map[string]*actorpool.ActorPool
}

func (pq *ActorPool) recalculateSelector() {
	if len(pq.peers) == 0 {
		pq.selector = selector{}
		return
	}

	peers := []string{}
	for name, pi := range pq.peers {
		if pi.state == dead || pi.optimisticState == dead {
			continue
		}
		if name == "" {
			panic("empty peer name")
		}
		peers = append(peers, name)
	}
	sort.Strings(peers)

	pq.selector.peers = peers
}

func (tr *Tracker) AddPool(pool *actorpool.ActorPool) error {
	_, ok := tr.pools[pool.ActorType()]
	if ok {
		return ErrActorTypeAlreadyRegistered
	}
	tr.pools[pool.ActorType()] = pool
	return nil
}

// Live peer.
func (tr *Tracker) Live(peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}
	pi.state = live
	pi.optimisticState = live

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

// OptimisticallyLive until an event marks the peer dead.
// Currently this has no affect on scheduling.
func (tr *Tracker) OptimisticallyLive(peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}
	pi.optimisticState = live

	// In the optimistic "live" case don't actually
	// recalculate the selector anything.
}

// Dead peer.
func (tr *Tracker) Dead(peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}
	pi.state = dead
	pi.optimisticState = dead

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

// OptimisticallyDead until a real event marks the peer alive again.
// Making the peer optimistically dead will remove it from any
// scheduling, in other words, it will never be returned as a
// peer from MinAssigned when marked optimistically dead.
func (tr *Tracker) OptimisticallyDead(peer string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if !isValidName(peer) {
		return
	}

	pi, ok := pq.peers[peer]
	if !ok {
		pi = newPeerInfo(peer)
		pq.peers[peer] = pi
	}
	pi.optimisticState = dead

	// Recalculate which peers are next
	// in the various selection schemes.
	pq.recalculateSelector()
}

func (tr *Tracker) Pools() map[string]*actorpool.ActorPool {
	tmpcp := map[string]*actorpool.ActorPool{}
	for key, pool := range tr.pools {
		tmpcp[key] = pool
	}
	return tmpcp
}

func (tr *Tracker) PoolByType(typ string) (*actorpool.ActorPool, error) {
	pool, ok := tr.pools[typ]
	if !ok {
		return nil, ErrUnknownActorType
	}
	return pool, nil
}

//returns a struct that represents all the peer queue's internal states used for logging and debugging
// relocation issues.
func (tr *Tracker) Status() *ClusterStatus {
	clusterState := map[string]*actorpool.PeersStatus{}
	for actortype, pool := range tr.pools {
		clusterState[actortype] = pool.Status()
	}
	return &ClusterStatus{
		ClusterState: clusterState,
	}
}

/*
func (tr *Tracker) BestPeer(def *grid.ActorStart) (string, error) {
	pool, ok := tr.pools[def.Type]
	if !ok {
		return "", ErrUnknownActorType
	}
	return pool.ByHash(def.Name)
}

func (tr *Tracker) Missing() []*grid.ActorStart {
	var all []*grid.ActorStart
	for _, pool := range tr.pools {
		defs := pool.Missing()
		all = append(all, defs...)
	}
	return all
}
*/
