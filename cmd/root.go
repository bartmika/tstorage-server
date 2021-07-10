package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tstorage-server",
	Short: "Time-series data storage gRPC server",
	Long: `The purpose of this application is to provide a local
	       database tailored for fast time-series data storage and be
		   accessible with remote procedure calls (gRPC).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
