package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the current version that this server is on.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tstorage-server v1.0")
	},
}
