package hash

import (
	"testing"

	"github.com/epsniff/gui/src/server/scheduler/actor/pool"
	"github.com/epsniff/gui/src/server/scheduler/tracker/clusterstate"
)

func TestHashPlacement_BestPeer(t *testing.T) {
	ps1 := clusterstate.New()
	ps1.Live("p1")
	ps1.Live("p2")
	ps1.Live("p3")
	ap1 := pool.New(ps1)

	map1 := map[string]*pool.ActorPool{"ap1": ap1}
	type args struct {
		actorName  string
		peersState clusterstate.PeersState
		pool       map[string]*pool.ActorPool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic1", args{"wrkr01", ps1, map1}, "p2", false},
		{"basic2", args{"wrkr02", ps1, map1}, "p2", false},
		{"basic3", args{"wrkr03", ps1, map1}, "p1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			got, err := h.BestPeer(tt.args.actorName, tt.args.peersState, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPlacement.BestPeer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashPlacement.BestPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPlacement_Relocate(t *testing.T) {
	t.Skip("not implemented")
	//	type args struct {
	//		peersState clusterstate.PeersState
	//		pools      map[string]*pool.ActorPool
	//	}
	//	tests := []struct {
	//		name    string
	//		h       *HashPlacement
	//		args    args
	//		want    *plan.RelocationPlan
	//		wantErr bool
	//	}{
	//		// TODO: Add test cases.
	//	}
	//	for _, tt := range tests {
	//		t.Run(tt.name, func(t *testing.T) {
	//			h := &HashPlacement{}
	//			got, err := h.Relocate(tt.args.peersState, tt.args.pools)
	//			if (err != nil) != tt.wantErr {
	//				t.Errorf("HashPlacement.Relocate() error = %v, wantErr %v", err, tt.wantErr)
	//				return
	//			}
	//			if !reflect.DeepEqual(got, tt.want) {
	//				t.Errorf("HashPlacement.Relocate() = %v, want %v", got, tt.want)
	//			}
	//		})
	//	}
}
