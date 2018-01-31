package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/araddon/gou"
	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/epsniff/spider/src/server/client"
	"github.com/epsniff/spider/src/server/definition"
	"github.com/epsniff/spider/src/server/name"
	"github.com/epsniff/spider/src/telemetry"
	"github.com/lytics/grid"
)

func RunServer(ctx context.Context, namespace, tcp, hostname string, port int, leaseDuration time.Duration, etcd *etcdv3.Client) error {

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
			return client.New(namespace, etcd)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create grid definition: %v", err)
	}
	server.RegisterDef(name.LeaderActor, def.MakeLeader)

	// Setup the gRPC grid server
	lis, err := net.Listen(tcp, fmt.Sprintf("%v:%d", hostname, port))
	if err != nil {
		return fmt.Errorf("failed listening on address and port: %s:%d: %v", hostname, port, err)
	}
	errout := make(chan error)
	go func() {
		defer close(errout)
		errout <- server.Serve(lis)
	}()

	// Wait for a shutdown signal or until we encounter an error
	select {
	case <-ctx.Done():
		telemetry.Logger.Infof("grid received stop signal")
		server.Stop()
		err = <-errout
		if err != nil {
			telemetry.Logger.Warnf("received unexpected grid shutdown error: %v", err)
		}
	case err = <-errout:
		telemetry.Logger.Warnf("received unexpected grid shutdown error: %v", err)
	}
	gou.Infof("grid shutdown complete")
	return err

}
