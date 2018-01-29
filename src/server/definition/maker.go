package definition

import (
	"github.com/epsniff/spider/src/server/actors/leader"
	"github.com/lytics/grid"
)

type Maker struct {
	//TODO add injection context
}

func New() (*Maker, error) {
	return &Maker{}, nil
}

func (m *Maker) MakeLeader(_ []byte) (grid.Actor, error) {
	cfg := &leader.Cfg{}
	return leader.New(cfg)
}
