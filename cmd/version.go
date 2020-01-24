package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xortim/peruse/conf"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  `Show version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(conf.GitVersion)
		},
	}
	return cmd
}
