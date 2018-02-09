package actorregistry

import (
	"github.com/epsniff/gui/src/server/actors/leader"
	"github.com/lytics/grid"
)

type Maker struct {
	clientMaker func() (*grid.Client, error)
}

func New(clientmaker func() (*grid.Client, error)) (*Maker, error) {
	return &Maker{clientmaker}, nil
}

func (m *Maker) MakeLeader(_ []byte) (grid.Actor, error) {
	cfg := &leader.Cfg{}
	gclient, err := m.clientMaker()
	if err != nil {
		return nil, err
	}
	return leader.New(gclient, cfg), nil
}
