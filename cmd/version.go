package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Set via -ldflags at build time; see Makefile.
var (
	version = "0.1.1"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("etoro %s (commit %s, built %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("etoro {{.Version}}\n")
}
