package actorpool

import "math"

type PeersStatus struct {
	AveActorsPerPeer     int                    `json:"ave_actors_peer"`
	Peers                map[string]*PeerState  `json:"peers"`
	Registered           map[string]*ActorState `json:"registered"`
	OptimisticRegistered map[string]*ActorState `json:"optimistic_registered"`
}

//PeerState is an reporting only changes to it don't effect internal state.
// This struct is a duplicate of peerInfo, but is exposed outside the package for reporting
// while maintaining the encapsulation of the internal state.
type PeerState struct {
	Name                 string   `json:"name"`           // Name of peer.
	State                bool     `json:"alive"`          // is the peer alive or dead.
	OptimisticState      bool     `json:"opt_alive"`      // is the peer optimisticly alive or dead.
	OptimisticRegistered []string `json:"actors_opt_reg"` // Possible actors on this peer.
	Registered           []string `json:"actors_reg"`     // Actors on this peer.
	Actors               int      `json:"actors_cnt"`     // count of actors on this peer.
	Burden               int      `json:"burden"`         // how much above or below the ave actors per peer this peer is.
}

//ActorStatus is an reporting only changes to it don't effect internal state.
type ActorState struct {
	Peer              string `json:"peer"`      // Name of peer.
	IsAlive           bool   `json:"alive"`     // Live or dead.
	IsOptimisticAlive bool   `json:"opt_alive"` // Live or dead.
}

//returns a struct that represents this peer queue's internal state.  Used for loggging.
func (pq *ActorPool) Status() *PeersStatus {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	//TODO add states like peer.burden, peer.actorCnt, pq.NumbActor
	alive := len(pq.registered)
	if alive <= 0 {
		return nil
	}

	pcnt := 0
	for _, p := range pq.peers {
		if p.state == live {
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
	for p, pi := range pq.peers {
		ps := &PeerState{
			Name:            pi.name,
			State:           pi.state,
			OptimisticState: pi.optimisticState,
		}
		reg := []string{}
		for actor, _ := range pi.registered {
			reg = append(reg, actor)
		}
		ps.Registered = reg
		optreg := []string{}
		for actor, _ := range pi.registered {
			optreg = append(optreg, actor)
		}
		ps.OptimisticRegistered = optreg
		ps.Actors = pi.NumActors()
		ps.Burden = len(pi.registered) - avePerPeer
		peers[p] = ps
	}

	registered := map[string]*ActorState{}
	for p, pi := range pq.registered {
		ps := &ActorState{
			Peer:              pi.name,
			IsAlive:           pi.state,
			IsOptimisticAlive: pi.optimisticState,
		}
		registered[p] = ps
	}

	optimisticRegistered := map[string]*ActorState{}
	for p, pi := range pq.optimisticRegistered {
		ps := &ActorState{
			Peer:              pi.name,
			IsAlive:           pi.state,
			IsOptimisticAlive: pi.optimisticState,
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
