package peertracker

import (
	"time"

	"github.com/lytics/grid"
	"github.com/lytics/retry"
)

PeerMonitorActorPrefix     = "peer-monitor"

// PeerMonitorActor name.
func PeerMonitorActor(peer string) string {
	return fmt.Sprintf("%v-%v", PeerMonitorActorPrefix, peer)
}

func StartPeerMonitor(client *grid.Client, peer string) error {
	def := grid.NewActorStart(PeerMonitorActor(peer))
	def.Type = PeerMonitorActorPrefix
	def.Data = []byte(peer)

	var err error
	retry.X(3, 5*time.Second, func() bool {
		_, err = sm.client.Request(actorStartTimeout, peer, def)
		return err != nil
	})
	return err
}

