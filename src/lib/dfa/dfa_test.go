package dfa

import (
	"testing"
)

func TestNew(t *testing.T) {
	dfa := New()

	notStarted := State(0)
	alive := State(1)
	paused := State(2)
	exiting := State(2)
	dead := State(3)

	dfa.SetTransition(notStarted, "running", alive)
	dfa.SetTransition(alive, "ping", alive)
	dfa.SetTransition(alive, "pause", paused)
	dfa.SetTransition(paused, "unpause", alive)
	dfa.SetTransition(alive, "exit-signal", exiting)
	dfa.SetTransition(exiting, "exited", dead)

	dfa.SetStartState(notStarted)

	state, err := dfa.Input("running")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != alive {
		t.Fatalf("unexpected state:%v", state)
	}

	_, err = dfa.Input("foobar")
	if err != ErrInvailedInput {
		t.Fatalf("we expected an ErrInvailedInput error, got:%v", err)
	}

	state, err = dfa.Input("ping")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != alive {
		t.Fatalf("unexpected state:%v", state)
	}

	state, err = dfa.Input("pause")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != paused {
		t.Fatalf("unexpected state:%v", state)
	}

	state, err = dfa.Input("exit-signal")
	if err != ErrInvailedInput {
		t.Fatalf("we expected an ErrInvailedInput error, got:%v", err)
	}

	state, err = dfa.Input("unpause")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != alive {
		t.Fatalf("unexpected state:%v", state)
	}

	state, err = dfa.Input("exit-signal")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != exiting {
		t.Fatalf("unexpected state:%v", state)
	}

	state, err = dfa.Input("exited")
	if err != nil {
		t.Fatalf("unexpected error:%v", err)
	}
	if state != dead {
		t.Fatalf("unexpected state:%v", state)
	}
}
