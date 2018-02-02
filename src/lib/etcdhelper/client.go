package etcdhelper

import etcdv3 "github.com/coreos/etcd/clientv3"

func NewEdtcClient(servers []string) (*etcdv3.Client, error) {
	return etcdv3.New(etcdv3.Config{Endpoints: servers})
}
