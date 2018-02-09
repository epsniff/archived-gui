package main

import (
	"log"
	"os"

	"github.com/epsniff/gui/cmds"
	"github.com/epsniff/gui/cmds/runserver"
	"github.com/epsniff/gui/src/lib/logging"
	"github.com/epsniff/gui/src/lib/logging/loggergou"
	"github.com/spf13/cobra"
)

func main() {
	logging.Logger = loggergou.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), "debug")

	root := &cobra.Command{
		Use:   "gui",
		Short: "gui command line tools and server services",
	}

	root.AddCommand(cmds.VersionCmd)
	root.AddCommand(runserver.ServerCmd)

	root.Execute()
}
