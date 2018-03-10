package peerstate

import "testing"

func TestPeerOptimisticallyLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}
	ps := New()

	for p := range peers {
		ps.OptimisticallyLive(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Dead {
			t.Fatal("expected live peer")
		}
		if optimisticState != Live {
			t.Fatal("expected live peer")
		}
	}
}

func TestPeerLive(t *testing.T) {
	t.Parallel()
	peers := map[string]bool{
		"peer0": true,
		"peer1": true,
	}

	ps := New()

	for p := range peers {
		ps.Live(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Live {
			t.Fatal("expected live peer")
		}
		if optimisticState != Live {
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

	ps := New()
	for p := range peers {
		ps.OptimisticallyLive(p)
	}
	for p := range peers {
		ps.Live(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Live {
			t.Fatal("expected live peer")
		}
		if optimisticState != Live {
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
	ps := New()

	for p := range peers {
		ps.Live(p)
	}
	for p := range peers {
		ps.OptimisticallyDead(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Live {
			t.Fatal("expected live peer")
		}
		if optimisticState != Dead {
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
	ps := New()

	for p := range peers {
		ps.Live(p)
	}
	for p := range peers {
		ps.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Dead {
			t.Fatal("expected live peer")
		}
		if optimisticState != Dead {
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
	ps := New()

	for p := range peers {
		ps.Live(p)
	}
	for p := range peers {
		ps.OptimisticallyDead(p)
	}
	for p := range peers {
		ps.Dead(p)
	}
	for p := range peers {
		state, optimisticState := ps.State(p)
		if state != Dead {
			t.Fatal("expected dead peer")
		}
		if optimisticState != Dead {
			t.Fatal("expected dead peer")
		}
	}
}

func TestDeadPeerThatWasNeverLive(t *testing.T) {
	t.Parallel()
	ps := New()

	ps.Dead("peer0")
	state, optimisticState := ps.State("peer0")
	if state != Dead {
		t.Fatal("expected dead peer")
	}
	if optimisticState != Dead {
		t.Fatal("expected dead peer")
	}
}
