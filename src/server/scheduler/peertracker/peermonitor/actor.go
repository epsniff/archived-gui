package peermonitor

import (
	"context"
	"time"

	"github.com/araddon/gou"
	"github.com/epsniff/spider/src/lib/logging"
	"github.com/epsniff/spider/src/server/name"
	"github.com/lytics/dfa"
	"github.com/lytics/grid"
)

var (
	// States
	Running   = dfa.State("running")
	Finishing = dfa.State("finishing")
	Exiting   = dfa.State("exiting")
	// Letters
	Exit    = dfa.Letter("exit")
	Failure = dfa.Letter("failure")
)

func defineDFA(a *Actor) *dfa.DFA {
	d := dfa.New()
	d.SetStartState(Running)
	d.SetTerminalStates(Exiting)

	d.SetTransition(Running, Exit, Finishing, a.Finishing)
	d.SetTransition(Running, Failure, Exiting, a.Exiting)

	d.SetTransition(Finishing, Exit, Exiting, a.Exiting)
	d.SetTransition(Finishing, Failure, Exiting, a.Exiting)

	return d
}

func New(client *grid.Client, peer string) (*Actor, error) {
	return &Actor{
		peer:   peer,
		client: client,
	}, nil
}

type Actor struct {
	id     string
	ctx    context.Context
	peer   string
	client *grid.Client
}

func (a *Actor) String() string {
	return a.id
}

func (a *Actor) Act(ctx context.Context) {
	a.ctx = ctx

	id, err := grid.ContextActorName(ctx)
	if err != nil {
		gou.Errorf("%v: error: %v", a, err)
		return
	}
	a.id = id

	a.RunDFA()
}

func (a *Actor) RunDFA() (dfa.State, bool) {
	d := defineDFA(a)

	d.SetTransitionLogger(func(state dfa.State) {
		logging.Logger.Infof("%v: switched to state: %v", a, state)
	})

	return d.Run(a.Running)
}

func (a *Actor) Running() dfa.Letter {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return Exit
		case <-timer.C:
			err := a.sendStopping(false)
			if err != nil {
				logging.Logger.Errorf("%v: failed to inform leader of pod status: %v", a, err)
			}
			timer.Reset(20 * time.Second)
		}
	}
}

const FinishingTickInterval = 5 * time.Second

func (a *Actor) Finishing() dfa.Letter {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-time.After(FinishingTickInterval):
			return Exit
		case <-timer.C:
			err := a.sendStopping(true)
			if err != nil {
				gou.Errorf("%v: failed to inform leader of pod status: %v", a, err)
			}
			timer.Reset(20 * time.Second)
		}
	}
}

func (a *Actor) Exiting() {}

func (a *Actor) sendStopping(stopping bool) error {
	_, err := a.client.Request(10*time.Second, name.LeaderActor, &msgs.PeerStatusMsg{
		Peer:     a.peer,
		Stopping: stopping,
	})
	return err
}
