package server

import (
	"fmt"
	"time"

	"github.com/epsniff/spider/src/server/definition"
	"github.com/epsniff/spider/src/telemetry"
	"github.com/lytics/grid"
)

func InitServer(namespace string, leaseDuration time.Duration) error {

	// Create a grid server configuration.
	cfg := grid.ServerCfg{
		Namespace:     namespace,
		Logger:        telemetry.Logger,
		LeaseDuration: leaseDuration,
	}

	// Register actor definitions.
	def, err := definition.New()
	if err != nil {
		return fmt.Errorf("failed to create grid definition: %v", err)
	}
}
