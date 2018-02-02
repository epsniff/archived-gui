package runserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epsniff/spider/src/lib/etcdhelper"
	"github.com/epsniff/spider/src/lib/logging"
	"github.com/epsniff/spider/src/lib/nethelper"
	"github.com/epsniff/spider/src/server"
	"github.com/lytics/grid"
	"github.com/spf13/cobra"
)

var (
	ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "startup a new node in the cluster",
		Run:   runServer,
	}
)

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stop via signal.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		logging.Logger.Infof("Shutting down...")
		cancel()
	}()

	// TODO startup promethious metrics endpoint

	hostIp, err := nethelper.BindableIP()
	if err != nil {
		logging.Logger.Errorf("error selecting a bindable IP address: err:%v", err)
	}
	const leaseDuration = 2 * time.Minute
	const namespace = "default_namespace"
	const port = 5503
	var etcdServers = []string{"localhost:2379"}
	etcdv3, err := etcdhelper.NewEdtcClient(etcdServers)
	if err != nil {
		logging.Logger.Errorf("error creating an etcd client with servers[%v]: err:%v", etcdServers, err)
	}

	for i := 0; i < 2; i++ {
		err := server.RunServer(ctx, namespace, hostIp, port, leaseDuration, etcdv3)
		switch {
		case err == nil || err == context.Canceled:
			logging.Logger.Infof("server shutdown complete")
			return
		case err == grid.ErrAlreadyRegistered && i == 0:
			// Wait for previous lease to die.
			logging.Logger.Infof("grid returned `%v`, waiting for previous lease to expire: %v", err, leaseDuration)
			select {
			case <-time.After(leaseDuration + time.Second):
			case <-ctx.Done():
				return
			}
		default:
			logging.Logger.Errorf("the server returned an unexpected error: %v", err)
			return
		}
	}
}
