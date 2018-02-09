package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/araddon/gou"
	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/epsniff/gui/src/lib/gridclient"
	"github.com/epsniff/gui/src/lib/logging"
	"github.com/epsniff/gui/src/server/actorregistry"
	"github.com/epsniff/gui/src/server/name"
	"github.com/lytics/grid"
)

func RunServer(ctx context.Context, namespace, bindAddress string, leaseDuration time.Duration, etcd *etcdv3.Client) error {

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
	actorRegistry, err := actorregistry.New(
		func() (*grid.Client, error) {
			return gridclient.New(namespace, etcd)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create grid definition: %v", err)
	}
	server.RegisterDef(name.LeaderActor, actorRegistry.MakeLeader)

	// Setup the gRPC grid server

	lis, err := createAddressListner("tcp", bindAddress)
	if err != nil {
		return fmt.Errorf("failed listening on:%v err:%v", bindAddress, err)
	}

	gou.Infof("starting server: bind_address:`%v`", bindAddress)
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

	return err

}

// createAddressListner takes an address as ip:port,
// validates it and then creates a net.Listner for it.
func createAddressListner(netType, bindAddress string) (net.Listener, error) {
	lis, err := net.Listen(netType, bindAddress)
	if err != nil {
		return nil, fmt.Errorf("listen error: err:%v", err)
	}

	addr := lis.Addr()
	switch addr := addr.(type) {
	default:
		return nil, fmt.Errorf("unsupported address type: %T", addr)
	case *net.TCPAddr:
		if addr.IP.IsUnspecified() {
			return nil, fmt.Errorf("ip not specified, grid doesn't support unspecified addresses: %v", bindAddress)
		}
		return lis, nil
	}
}
