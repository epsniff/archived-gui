package actorpool

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
