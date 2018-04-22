package tracker

import (
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/actor/placement"
	actorpool "github.com/epsniff/gui/src/server/scheduler/actor/pool"
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
	"github.com/epsniff/gui/src/server/scheduler/types"
	"github.com/lytics/grid"
)

type poolname string

func New() *Tracker {
	return &Tracker{
		pools:     map[string]*actorpool.ActorPool{},
		placers:   map[string]placement.Placement{},
		peerState: clusterstate.New(),
	}
}

type Tracker struct {
	mu        sync.RWMutex
	pools     map[string]*actorpool.ActorPool
	placers   map[string]placement.Placement
	peerState clusterstate.PeersState
}

func (tr *Tracker) BestPeer(poolName string, def *grid.ActorStart) (string, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	placer, ok := tr.placers[poolName]
	if !ok {
		return "", types.ErrUnknownPoolName
	}
	return placer.BestPeer(def.Name, tr.peerState, tr.pools)
}

func (tr *Tracker) CreateActorPool(poolName string, placer placement.Placement) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	_, ok := tr.pools[poolName]
	if ok {
		return types.ErrActorPoolAlreadyRegistered
	}

	ap := actorpool.New(tr.peerState)
	tr.pools[poolName] = ap
	tr.placers[poolName] = placer

	return nil
}

func (tr *Tracker) Live(peer string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return tr.peerState.Live(peer)
}

func (tr *Tracker) Dead(peer string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	return tr.peerState.Dead(peer)
}

func (tr *Tracker) OptimisticallyLive(peer string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return tr.peerState.OptimisticallyLive(peer)
}

func (tr *Tracker) OptimisticallyDead(peer string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	return tr.peerState.OptimisticallyDead(peer)
}

func (tr *Tracker) Register(poolName string, isRequired bool, actor *grid.ActorStart, peer string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	pool, ok := tr.pools[poolName]
	if !ok {
		return types.ErrUnknownPoolName
	}

	return pool.Register(isRequired, actor, peer)

}

func (tr *Tracker) Pools() map[string]*actorpool.ActorPool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tmpCp := map[string]*actorpool.ActorPool{}
	for key, pool := range tr.pools {
		tmpCp[key] = pool
	}
	return tmpCp
}

func (tr *Tracker) PoolBy(name string) (*actorpool.ActorPool, error) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	pool, ok := tr.pools[name]
	if !ok {
		return nil, types.ErrUnknownPoolName
	}
	return pool, nil
}

func (tr *Tracker) Missing() []*grid.ActorStart {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	var all []*grid.ActorStart
	for _, pool := range tr.pools {
		defs := pool.Missing()
		all = append(all, defs...)
	}
	return all
}
