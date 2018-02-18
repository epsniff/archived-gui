package peertracker

import (
	"errors"
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/actorpool"
)

var (
	ErrActorTypeAlreadyRegistered = errors.New("actor type already registered")
	ErrUnknownActorType           = errors.New("unknown actor type")
)

func New() *Tracker {
	return &Tracker{
		pools: map[string]*actorpool.ActorPool{},
	}
}

type Tracker struct {
	mu    sync.Mutex
	pools map[string]*actorpool.ActorPool
}

func (tr *Tracker) AddPool(pool *actorpool.ActorPool) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	_, ok := tr.pools[pool.ActorType()]
	if ok {
		return ErrActorTypeAlreadyRegistered
	}
	tr.pools[pool.ActorType()] = pool
	return nil
}

func (tr *Tracker) Live(peer string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for _, pool := range tr.pools {
		pool.Live(peer)
	}
}

func (tr *Tracker) Dead(peer string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for _, pool := range tr.pools {
		pool.Dead(peer)
	}
}

func (tr *Tracker) OptimisticallyLive(peer string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for _, pool := range tr.pools {
		pool.OptimisticallyLive(peer)
	}
}

func (tr *Tracker) OptimisticallyDead(peer string) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for _, pool := range tr.pools {
		pool.OptimisticallyDead(peer)
	}
}

func (tr *Tracker) Pools() map[string]*actorpool.ActorPool {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	tmpCp := map[string]*actorpool.ActorPool{}
	for key, pool := range tr.pools {
		tmpCp[key] = pool
	}
	return tmpCp
}

func (tr *Tracker) PoolByType(typ string) (*actorpool.ActorPool, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	pool, ok := tr.pools[typ]
	if !ok {
		return nil, ErrUnknownActorType
	}
	return pool, nil
}

//returns a struct that represents all the peer queue's internal states used for logging and debugging
// relocation issues.
func (tr *Tracker) Status() *ClusterStatus {
	tr.mu.Lock()
	defer tr.mu.Unlock()

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
