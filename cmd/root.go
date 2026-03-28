package cmd

import (
	"fmt"
	"os"

	"github.com/marianopa-tr/etoro-cli/internal/config"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	flagOutput  string
	flagDemo    bool
	flagYes     bool
	flagTimeout string
)

var rootCmd = &cobra.Command{
	Use:   "etoro",
	Short: "eToro CLI — trade, invest, and copy from your terminal",
	Long: `eToro CLI — trade, invest, and copy from your terminal.

Browse markets, manage your portfolio, place trades, and interact with
eToro's social trading platform — from a terminal or as a JSON API for
scripts and AI agents.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch flagOutput {
		case "json":
			output.SetFormat(output.JSON)
		default:
			output.SetFormat(output.Table)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if output.GetFormat() == output.JSON {
			fmt.Fprintf(os.Stdout, `{"error":%q}`+"\n", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "", `output format: "table" (default) or "json"`)
	rootCmd.PersistentFlags().BoolVar(&flagDemo, "demo", false, "use demo/virtual account for trading")
	rootCmd.PersistentFlags().BoolVar(&flagYes, "yes", false, "skip confirmation prompts")
	rootCmd.PersistentFlags().StringVar(&flagTimeout, "timeout", "", "request timeout (e.g. 30s, 1m)")
}

func mustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		output.Errorf("failed to load config: %s", err)
		os.Exit(1)
	}
	return cfg
}
