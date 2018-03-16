package plan

import "github.com/lytics/grid"

func New() *RelocationPlan {
	return &RelocationPlan{[]*Relocation{}}
}

type Relocation struct {
	PoolName  string
	ActorName string
	Def       *grid.ActorStart
}

type RelocationPlan struct {
	Relocations []*Relocation
}
