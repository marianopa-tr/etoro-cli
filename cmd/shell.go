package cmd

import (
	"reflect"
	"unsafe"

	"github.com/etoro/etoro-cli/internal/api"
	"github.com/etoro/etoro-cli/internal/resolver"
	"github.com/etoro/etoro-cli/internal/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive REPL with history and tab-completion",
	Long: `Start an interactive eToro CLI shell with command history,
tab-completion for commands and symbols, and a persistent session.

Examples:
  etoro shell`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()

		if cfg.Auth.APIKey == "" || cfg.Auth.UserKey == "" {
			setupCmd.RunE(setupCmd, nil)
			cfg = mustLoadConfig()
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		commands := []string{
			"search", "instruments get", "quote",
			"portfolio summary", "portfolio positions", "portfolio orders", "portfolio history",
			"trade open", "trade close", "trade limit",
			"orders list", "orders cancel", "orders cancel-all",
			"watchlist list", "watchlist create", "watchlist add", "watchlist remove",
			"copy discover", "copy performance", "copy copiers",
			"feed list", "feed post",
			"pi copiers", "pi get", "pi gain",
			"status", "auth status", "setup",
		}

		symbols := res.CachedSymbols()

		sh := shell.New(func(shellArgs []string) error {
			resetCommandTree(rootCmd)
			rootCmd.SetArgs(shellArgs)
			return rootCmd.Execute()
		}, commands, symbols)

		return sh.Run()
	},
}

// resetCommandTree resets Cobra/pflag state so commands can be
// re-executed in a REPL loop. pflag marks its FlagSet as "parsed"
// after the first Parse() call and skips all subsequent parses,
// which breaks flag and positional-arg handling on the second run.
func resetCommandTree(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	resetFlagSetParsed(cmd.Flags())

	for _, sub := range cmd.Commands() {
		resetCommandTree(sub)
	}
}

func resetFlagSetParsed(fs *pflag.FlagSet) {
	fv := reflect.ValueOf(fs).Elem()
	if field := fv.FieldByName("parsed"); field.IsValid() {
		ptr := unsafe.Pointer(field.UnsafeAddr())
		*(*bool)(ptr) = false
	}
	if field := fv.FieldByName("args"); field.IsValid() {
		ptr := unsafe.Pointer(field.UnsafeAddr())
		*(*[]string)(ptr) = nil
	}
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
