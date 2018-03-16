package scheduler

import (
	"context"
	"time"

	"github.com/epsniff/gui/src/lib/logging"
	"github.com/epsniff/gui/src/server/scheduler/tracker"
	"github.com/lytics/grid"
)

const actorStartTimeout = 10 * time.Second

type Scheduler struct {
	ctx     context.Context
	tracker *tracker.Tracker
	client  *grid.Client
}

func New(ctx context.Context, client *grid.Client) *Scheduler {
	s := &Scheduler{
		ctx:     ctx,
		tracker: tracker.New(),
		client:  client,
	}

	return s
}

func (sm *Scheduler) Run() error {
	current, peers, err := sm.client.QueryWatch(sm.ctx, grid.Peers)
	if err != nil {
		logging.Logger.Errorf("%v: fatal error: %v", sm, err)
		return err
	}

	logging.Logger.Infof("%v: found %v current peers", sm, len(current))
	for _, c := range current {
		logging.Logger.Infof("%v: found existing peer: %v", sm, c.Peer())
		sm.tracker.Live(c.Peer())

		if err := tracker.StartPeerMonitor(sm.client, c.Peer()); err != nil {
			logging.Logger.Warnf("%v: failed to start peer monitor on: %v, error: %v", sm, c.Peer(), err)
		}
	}

	for {
		select {
		case <-sm.ctx.Done():
			return nil
		case e := <-peers:
			logging.Logger.Infof("%v: %v", sm, e)
			switch e.Type {
			case grid.WatchError:
				logging.Logger.Errorf("%v: fatal error: %v", sm, e.Err())
				return err
			case grid.EntityLost:
				sm.tracker.Dead(e.Peer())
			case grid.EntityFound:
				sm.tracker.Live(e.Peer())
				if err := tracker.StartPeerMonitor(sm.client, e.Peer()); err != nil {
					logging.Logger.Warnf("%v: failed to start peer monitor on: %v, error: %v", sm, e.Peer(), err)
				}
			}
		}
	}
}

func (sm *Scheduler) startActor(def *grid.ActorStart) error {
	pool := sm.tracker.PoolByType(def.GetType())

	peer, err := pool.BestPeer(def)
	if err != nil {
		return err
	}
	logging.Logger.Infof("%v: starting actor: %v, on peer: %v", sm, def.Name, peer)
	pool.OptimisticallyRegister(def.Name, peer)
	_, err = sm.client.Request(actorStartTimeout, peer, def)
	if err != nil {
		pool.OptimisticallyUnregister(def.Name)
		return err
	}
	return nil
}
