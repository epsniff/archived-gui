package gridclient

import (
	"runtime"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/epsniff/spider/src/lib/logging"
	"github.com/lytics/grid"
)

func New(namespace string, etcd *etcdv3.Client) (*grid.Client, error) {
	cfg := grid.ClientCfg{
		Logger:             logging.Logger,
		Namespace:          namespace,
		ConnectionsPerPeer: runtime.GOMAXPROCS(0) + 1,
	}
	client, err := grid.NewClient(etcd, cfg)
	if err != nil {
		return nil, err
	}
	logging.Logger.Debugf("grid client for namespace: %v", namespace)
	return client, nil
}
