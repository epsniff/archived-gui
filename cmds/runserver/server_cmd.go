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
	// flags and defaults
	address       = ":5503"
	leaseDuration = 2 * time.Minute
	namespace     = "default_namespace"
	etcdServers   = []string{"localhost:2379"}
)

func init() {
	ServerCmd.PersistentFlags().StringSliceVar(&etcdServers, "etcds", etcdServers, `which etcd servers to connect to`)
	ServerCmd.PersistentFlags().StringVar(&namespace, "namespace", namespace, `
		used to determine which namespace this grid should use. You'll need this if 
		your running mutiple grid servers on the same etcd cluster.`)
	ServerCmd.PersistentFlags().DurationVar(&leaseDuration, "lease_dur", leaseDuration, `
		used to determine how long a peer can be missing before it's considered down`)
	ServerCmd.PersistentFlags().StringVar(&address, "address", address, `
		The 'ip:port' to bind grid's internal gRPC server to, e.g. 127.0.0.1:5503.  
		If given in the form of ':5503' then we'll pick a non loopback address to bind 
		to, as grid doesn't support binding to all interfaces.`)
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Extract flags
	a, err := nethelper.ValidateAddress(address)
	if err != nil {
		logging.Logger.Errorf("address failed validation: address:%v err:%v", address, err)
	}
	address = a

	// Stop via signal.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		logging.Logger.Infof("Shutting down...")
		cancel()
	}()

	// TODO startup promethious metrics endpoint

	etcdv3, err := etcdhelper.NewEdtcClient(etcdServers)
	if err != nil {
		logging.Logger.Errorf("error creating an etcd client with servers[%v]: err:%v", etcdServers, err)
	}

	for i := 0; i < 2; i++ {
		err := server.RunServer(ctx, namespace, address, leaseDuration, etcdv3)
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
			logging.Logger.Errorf("%v", err)
			return
		}
	}
}
