package main

import (
	"log"
	"os"

	"github.com/epsniff/spider/cmds"
	"github.com/epsniff/spider/telemetry"
	"github.com/epsniff/spider/telemetry/loggergou"
	"github.com/spf13/cobra"
)

func main() {
	telemetry.Logger = loggergou.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), "debug")

	root := &cobra.Command{
		Use:   "spider",
		Short: "spider command line tools and server services",
	}

	root.AddCommand(cmds.VersionCmd)

	root.Execute()
}
