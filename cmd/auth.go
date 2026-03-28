package cmd

import (
	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long: `Check authentication status and manage API credentials.

Examples:
  etoro auth status`,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status and account info",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()

		hasAPIKey := cfg.Auth.APIKey != ""
		hasUserKey := cfg.Auth.UserKey != ""

		result := map[string]any{
			"apiKey":  hasAPIKey,
			"userKey": hasUserKey,
		}

		if hasAPIKey && hasUserKey {
			client := api.NewClient(cfg, flagDemo)
			_, err := client.GetPeople()
			if err != nil {
				result["authenticated"] = false
				result["error"] = err.Error()
			} else {
				result["authenticated"] = true
			}
		} else {
			result["authenticated"] = false
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(result)
			return nil
		}

		t := output.NewDetailTable()
		output.DetailRow(t, "API Key", output.FormatBool(hasAPIKey, "configured", "missing"))
		output.DetailRow(t, "User Key", output.FormatBool(hasUserKey, "configured", "missing"))

		if auth, ok := result["authenticated"].(bool); ok {
			if auth {
				output.DetailRow(t, "Status", output.Green("Authenticated"))
			} else {
				errMsg := "not authenticated"
				if e, ok := result["error"].(string); ok {
					errMsg = e
				}
				output.DetailRow(t, "Status", output.Red(errMsg))
			}
		}
		output.RenderTable(t)

		if !hasAPIKey || !hasUserKey {
			output.Infof("\nRun `etoro setup` to configure your API keys.")
		}
		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
