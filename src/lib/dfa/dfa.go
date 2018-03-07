package dfa

import (
	"fmt"
)

//
// https://en.wikipedia.org/wiki/Deterministic_finite_automaton
//

//
// finite set of states (Q)
//
type State int

func (s State) String() string {
	return fmt.Sprintf("%d", int(s))
}

//
// letters in the alphabet, which is a finite set of input symbols called the alphabet (Σ)
//
type Letter string

func (l Letter) String() string {
	return fmt.Sprintf("%s", string(l))
}

//
// transition function (δ : Q × Σ → Q)
//
type transition struct {
	state  State
	letter Letter
}

type DFA struct {
	startState  State
	state       State
	transitions map[transition]State
}

func New() *DFA {
	return &DFA{
		startState:  State(-1),
		state:       State(-1),
		transitions: map[transition]State{},
	}
}

func (dfa *DFA) SetStartState(start State) {
	dfa.startState = start
	dfa.state = start
}

func (dfa *DFA) SetTransition(from State, input Letter, to State) {
	_, ok := dfa.transitions[transition{from, input}]
	if ok {
		panic(fmt.Sprintf("transition for state:%v on letter:%v already set", from, input))
	}
	dfa.transitions[transition{from, input}] = to
}

func (dfa *DFA) printTransitions() {
	for trans, st := range dfa.transitions {
		fmt.Printf(" (at:%v on:%v) -- > %v\n", trans.state, trans.letter, st)
	}
}

var ErrInvailedInput = fmt.Errorf("invailed input for current state")

func (dfa *DFA) Input(input Letter) (State, error) {
	//fmt.Printf("INPUT (at:%v on:%v)\n", dfa.state, input)
	//dfa.printTransitions()
	state, ok := dfa.transitions[transition{dfa.state, input}]
	if !ok {
		return State(-1), ErrInvailedInput
	}
	dfa.state = state
	return state, nil
}
