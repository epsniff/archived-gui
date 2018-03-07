package peertracker

import (
	"testing"

	"github.com/epsniff/gui/src/server/scheduler/actorpool"
)

func TestPeerOptimisticallyLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	tr := New()
	tr.AddPool("pool1", actorpool.New(true))
	tr.AddPool("pool2", actorpool.New(true))
	for p := range peers {
		tr.OptimisticallyLive(p)
	}

	if err := tr.Register("pool1", "worker1", "peer1"); err != nil {
		t.Fatalf("err:%v", err)
	}

	pools := tr.Pools()
	if len(pools) != 2 {
		t.Fatalf("unexpected number of pools. %v", pools)
	}

	p1 := pools["pool1"]
	p2 := pools["pool2"]

	if !p1.IsRegistered("worker1") {
		t.Fatalf("we expected worker1 to be registered on pool1")
	}

	if p1.NumRegistered() != 1 {
		t.Fatalf("we expected pool1 to have one actor")
	}

	if p2.IsRegistered("worker1") || p2.NumRegistered() != 0 {
		t.Fatalf("we expected pool2 to have zero actors")
	}
}

/*
func TestPeerLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ap := New()
	ap.AddPool(actorpool.New(true))

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != live {
			t.Fatal("expected live peer")
		}
	}
}
func TestPeerOptimisticallyLiveToLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ap := New()
	ap.AddPool(actorpool.New(true))
	for p := range peers {
		ap.OptimisticallyLive(p)
	}
	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != live {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerOptimisticallyDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New()
	ap.AddPool(actorpool.New(true))

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.OptimisticallyDead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != live {
			t.Fatal("expected live peer")
		}
		if optimisticState != dead {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New()
	ap.AddPool(actorpool.New(true))

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != dead {
			t.Fatal("expected live peer")
		}
		if optimisticState != dead {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerOptimisticallyDeadToDead(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ap := New()
	ap.AddPool(actorpool.New(true))

	for p := range peers {
		ap.Live(p)
	}
	for p := range peers {
		ap.OptimisticallyDead(p)
	}
	for p := range peers {
		ap.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ap.peerState.State(p)
		if state != dead {
			t.Fatal("expected dead peer")
		}
		if optimisticState != dead {
			t.Fatal("expected dead peer")
		}
	}
}

func TestDeadPeerThatWasNeverLive(t *testing.T) {
	t.Parallel()
	ap := New()
	ap.AddPool(actorpool.New(true))

	ap.Dead("peer0")
	state, optimisticState := ap.peerState.State("peer0")
	if state != dead {
		t.Fatal("expected dead peer")
	}
	if optimisticState != dead {
		t.Fatal("expected dead peer")
	}
}
*/
