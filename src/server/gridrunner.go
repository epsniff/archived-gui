package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/araddon/gou"
	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/epsniff/spider/src/lib/gridclient"
	"github.com/epsniff/spider/src/lib/logging"
	"github.com/epsniff/spider/src/server/definition"
	"github.com/epsniff/spider/src/server/name"
	"github.com/lytics/grid"
	"github.com/lytics/grid/registry"
)

func RunServer(ctx context.Context, namespace, hostIp string, port int, leaseDuration time.Duration, etcd *etcdv3.Client) error {

	// Create a grid server configuration.
	cfg := grid.ServerCfg{
		Namespace:     namespace,
		LeaseDuration: leaseDuration,
	}

	// Create the server.
	server, err := grid.NewServer(etcd, cfg)
	if err != nil {
		return fmt.Errorf("failed to create grid server: %v", err)
	}

	// Register actor definitions.
	def, err := definition.New(
		func() (*grid.Client, error) {
			return gridclient.New(namespace, etcd)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create grid definition: %v", err)
	}
	server.RegisterDef(name.LeaderActor, def.MakeLeader)

	// Setup the gRPC grid server
	listenA := fmt.Sprintf("%v:%d", hostIp, port)
	lis, err := net.Listen("tcp", listenA)
	if err != nil {
		return fmt.Errorf("failed listening on:%v err:%v", listenA, err)
	}

	gou.Infof("starting server: bind_address:`%v`", listenA)
	errout := make(chan error)
	go func() {
		defer close(errout)
		errout <- server.Serve(lis)
	}()

	// Wait for a shutdown signal or until we encounter an error
	select {
	case <-ctx.Done():
		logging.Logger.Infof("grid received stop signal")
		server.Stop()
		err = <-errout
	case err = <-errout:
	}

	switch err {
	case nil:
	case registry.ErrUnspecifiedNetAddressIP:
		logging.Logger.Errorf("received bad address error from grid: address: `%v` error: %v", listenA, err)
	default:
		logging.Logger.Errorf("received unexpected grid error: %v", err)
	}

	return err

}
