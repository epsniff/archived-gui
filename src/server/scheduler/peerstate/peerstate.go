package peerstate

import (
	"errors"
	"sync"

	"github.com/epsniff/gui/src/server/scheduler/peerinfo"
)

const (
	Live = true
	Dead = false
)

var (
	ErrInvalidPeerName = errors.New("invalid peer name")
	ErrUnknownPeerName = errors.New("unknown peer name")
)

type PeersState interface {
	State(peer string) (state bool, optimisticstate bool)
	Live(peer string) error
	OptimisticallyLive(peer string) error
	Dead(peer string) error
	OptimisticallyDead(peer string) error
	Get(peer string) (*peerinfo.PeerInfo, error)
}

func New() PeersState {
	ps := &peersState{
		mu:    &sync.RWMutex{},
		peers: map[string]*peerinfo.PeerInfo{},
	}

	return ps
}

// peersState contains state of all the known and optimisiticaly started actors
// and which peers they have been started on.
type peersState struct {
	mu *sync.RWMutex

	peers map[string]*peerinfo.PeerInfo
}

func (ps *peersState) Get(peer string) (*peerinfo.PeerInfo, error) {
	if !isValidName(peer) {
		return nil, ErrInvalidPeerName
	}

	ps.mu.RLock()
	defer ps.mu.RUnlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return nil, ErrUnknownPeerName
	}
	return pi, nil
}

// State retrieves peer state.
func (ps *peersState) State(peer string) (state bool, optimisticstate bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		return false, false
	}
	state, optimisticstate = pi.State, pi.OptimisticState
	return state, optimisticstate
}

// Live peer.
func (ps *peersState) Live(peer string) error {
	if !isValidName(peer) {
		return ErrInvalidPeerName
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.State = Live
	pi.OptimisticState = Live
	return nil
}

// OptimisticallyLive until an event marks the peer dead.
// Currently this has no affect on scheduling.
func (ps *peersState) OptimisticallyLive(peer string) error {

	if !isValidName(peer) {
		return ErrInvalidPeerName
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.OptimisticState = Live
	return nil
}

// Dead peer.
func (ps *peersState) Dead(peer string) error {
	if !isValidName(peer) {
		return ErrInvalidPeerName
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.State = Dead
	pi.OptimisticState = Dead
	return nil
}

// OptimisticallyDead until a real event marks the peer alive again.
// Making the peer optimistically dead will remove it from any
// scheduling, in other words, it will never be returned as a
// peer from MinAssigned when marked optimistically dead.
func (ps *peersState) OptimisticallyDead(peer string) error {
	if !isValidName(peer) {
		return ErrInvalidPeerName
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	pi, ok := ps.peers[peer]
	if !ok {
		pi = peerinfo.New(peer)
		ps.peers[peer] = pi
	}
	pi.OptimisticState = Dead
	return nil
}

func isValidName(name string) bool {
	return name != ""
}
