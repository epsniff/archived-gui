package cmds

import (
	"fmt"

	"github.com/epsniff/spider/src"
	"github.com/spf13/cobra"
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of this binary",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version.Version)
		},
	}
)
