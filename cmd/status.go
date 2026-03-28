package cmd

import (
	"fmt"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check API connectivity and authentication status",
	Long: `Verify that your API keys are configured and the eToro API is reachable.

Examples:
  etoro status
  etoro status --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		client := api.NewClient(cfg, flagDemo)

		hasAPIKey := cfg.Auth.APIKey != ""
		hasUserKey := cfg.Auth.UserKey != ""

		result := map[string]any{
			"apiKey":  hasAPIKey,
			"userKey": hasUserKey,
			"demo":    flagDemo,
		}

		if hasAPIKey && hasUserKey {
			_, err := client.GetPortfolio()
			if err != nil {
				result["connected"] = false
				result["error"] = err.Error()
			} else {
				result["connected"] = true
			}
		} else {
			result["connected"] = false
			result["error"] = "missing credentials"
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(result)
			return nil
		}

		t := output.NewDetailTable()
		output.DetailRow(t, "API Key", output.FormatBool(hasAPIKey, "configured", "missing"))
		output.DetailRow(t, "User Key", output.FormatBool(hasUserKey, "configured", "missing"))
		output.DetailRow(t, "Mode", modeLabel())

		if connected, ok := result["connected"].(bool); ok && connected {
			output.DetailRow(t, "Connection", output.Green("OK"))
		} else {
			errMsg := "not tested"
			if e, ok := result["error"].(string); ok {
				errMsg = e
			}
			output.DetailRow(t, "Connection", output.Red(errMsg))
		}
		output.RenderTable(t)

		if !hasAPIKey || !hasUserKey {
			fmt.Println()
			output.Infof("Run `etoro setup` to configure your API keys.")
		}

		return nil
	},
}

func modeLabel() string {
	if flagDemo {
		return output.Yellow("Demo")
	}
	return "Real"
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
