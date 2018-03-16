package tracker

import (
	"errors"
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/actor/placement"
	actorpool "github.com/epsniff/gui/src/server/scheduler/actor/pool"
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
	"github.com/lytics/grid"
)

var (
	ErrActorPoolAlreadyRegistered = errors.New("actor pool name already registered")
	ErrUnknownPoolName            = errors.New("unknown actor pool name")
)

type poolname string

func New() *Tracker {
	return &Tracker{
		pools:     map[string]*actorpool.ActorPool{},
		peerState: clusterstate.New(),
	}
}

type Tracker struct {
	mu        sync.Mutex
	pools     map[string]*actorpool.ActorPool
	peerState clusterstate.PeersState
}

func (tr *Tracker) CreateActorPool(name string, placer placement.Placement) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	_, ok := tr.pools[name]
	if ok {
		return ErrActorPoolAlreadyRegistered
	}

	ap := actorpool.New(tr.peerState)
	tr.pools[name] = ap
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
		return ErrUnknownPoolName
	}

	return pool.Register(isRequired, actor, peer)

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

func (tr *Tracker) PoolBy(name string) (*actorpool.ActorPool, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	pool, ok := tr.pools[name]
	if !ok {
		return nil, ErrUnknownPoolName
	}
	return pool, nil
}

func (tr *Tracker) Missing() []*grid.ActorStart {
	var all []*grid.ActorStart
	for _, pool := range tr.pools {
		defs := pool.Missing()
		all = append(all, defs...)
	}
	return all
}

/*
func (tr *Tracker) BestPeer(def *grid.ActorStart) (string, error) {
	pool, ok := tr.pools[def.Type]
	if !ok {
		return "", ErrUnknownActorType
	}
	return pool.ByHash(def.Name)
}
*/
