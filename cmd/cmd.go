// Package cmd leverage spf13/cobra to create a CLI utility
package cmd

import (
	"os"
)

// Execute runs the root command
func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
