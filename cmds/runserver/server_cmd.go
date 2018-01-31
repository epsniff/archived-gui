package runserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epsniff/spider/src/server"
	"github.com/epsniff/spider/src/telemetry"
	"github.com/lytics/grid"
	"github.com/lytics/lio/src/lib/ipaddress"
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
		telemetry.Logger.Infof("Shutting down...")
		cancel()
	}()

	// TODO startup promethious metrics endpoint

	// Find the address to use. On G[K|C]E 'eth0' is the default.
	ip, err := ipaddress.FindInterface("eth0")
	if ip == nil {
		// If 'eth0' doesn't present an IP, try all interfaces.
		ip, err = ipaddress.Find()
	}
	tcp := "tcp"
	if ip.To16() != nil && ip.To4() == nil {
		tcp = "tcp6"
	}
	const leaseDuration = 2 * time.Minute

	for i := 0; i < 2; i++ {
		err := server.RunServer(ctx, namespace, tcp, hostname, port, leaseDuration, etcd_client)
		if err == grid.ErrAlreadyRegistered && i == 0 {
			// Wait for previous lease to die.
			telemetry.Logger.Infof("grid returned `%v`, waiting for previous lease to expire: %v", err, leaseDuration)
			select {
			case <-time.After(leaseDuration + time.Second):
			case <-ctx.Done():
				return
			}
		} else if err == context.Canceled {
			return
		} else {
			telemetry.Logger.Errorf("the server returned an unexpected error: %v", err)
			return
		}
	}
}
