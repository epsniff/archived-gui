package peertracker

import (
	"encoding/json"
	"fmt"

	"github.com/epsniff/gui/src/server/scheduler/actorpool"
)

type ClusterStatus struct {
	ClusterState map[string]*actorpool.PeersStatus `json:"cluster_state"`
}

func (cs *ClusterStatus) String() string {
	if cs == nil {
		return `{"cluster_state": null}`
	}
	if data, err := json.Marshal(cs); err != nil {
		return fmt.Sprintf(`{"cluster_state":"marshal_error=%v"}`, err)
	} else {
		return string(data)
	}
}

func (cs *ClusterStatus) Strings() []string {
	if cs == nil {
		return []string{`{"cluster_state": null}`}
	}
	res := []string{}
	for pqname, pqstatus := range cs.ClusterState {
		if data, err := json.Marshal(struct {
			PeerQueueName string                 `json:"peerqueue_name,omitempty"`
			State         *actorpool.PeersStatus `json:"state,omitempty"`
		}{
			PeerQueueName: pqname,
			State:         pqstatus,
		}); err != nil {
			return []string{fmt.Sprintf(`{"cluster_state":"marshal_error=%v"}`, err)}
		} else {
			res = append(res, string(data))
		}
	}
	return res
}
