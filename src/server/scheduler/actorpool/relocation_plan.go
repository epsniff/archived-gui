package actorpool

import (
	"bytes"
	"fmt"
	"sort"
)

func NewRelocationPlan(actorType string, total, average int) *RelocationPlan {
	return &RelocationPlan{
		ActorType:   actorType,
		Total:       total,
		Average:     average,
		Peers:       []string{},
		Count:       map[string]int{},
		Burden:      map[string]int{},
		Relocations: []string{}, //actor name. aka the mailbox
	}
}

type RelocationPlan struct {
	ActorType   string
	Total       int
	Average     int
	Peers       []string
	Count       map[string]int
	Burden      map[string]int
	Relocations []string //actor name. aka the mailbox
}

func (p *RelocationPlan) String() string {
	var peers []string
	for peer := range p.Count {
		peers = append(peers, peer)
	}
	sort.Strings(peers)

	var count bytes.Buffer
	for i, peer := range peers {
		count.WriteString(fmt.Sprintf("%v=%v", peer, p.Count[peer]))
		if i+1 < len(peers) {
			count.WriteString(", ")
		}
	}
	var burden bytes.Buffer
	for i, peer := range peers {
		burden.WriteString(fmt.Sprintf("%v=%v", peer, p.Burden[peer]))
		if i+1 < len(peers) {
			burden.WriteString(", ")
		}
	}

	return fmt.Sprintf("actor-type: %v; nr-peers: %v; peers: %v; nr-actors: %v; target-per-peer: %v; count: %v; burden: %v",
		p.ActorType,
		len(p.Peers),
		p.Peers,
		p.Total,
		p.Average,
		count.String(),
		burden.String())
}
