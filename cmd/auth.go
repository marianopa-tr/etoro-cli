package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/config"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long: `Log in, log out, and check authentication status.

Examples:
  etoro auth login --api-key <key> --user-key <key>
  etoro auth login
  etoro auth logout
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
			identity, err := client.GetIdentity()
			if err != nil {
				result["authenticated"] = false
				result["error"] = err.Error()
			} else {
				result["authenticated"] = true
				result["gcid"] = identity.GCID
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
				if gcid, ok := result["gcid"].(int); ok {
					output.DetailRow(t, "GCID", gcid)
				}
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

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your eToro API keys",
	Long: `Store your API keys and verify they work.
Pass keys via flags for non-interactive use, or omit them to be prompted.

Get your keys from: https://www.etoro.com/settings/trade

Examples:
  etoro auth login --api-key <key> --user-key <key>
  etoro auth login
  ETORO_PUBLIC_KEY=... ETORO_USER_KEY=... etoro auth login`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		apiKey, _ := cmd.Flags().GetString("api-key")
		userKey, _ := cmd.Flags().GetString("user-key")

		if apiKey == "" {
			apiKey = promptSecret("API Key (ETORO_PUBLIC_KEY)", cfg.Auth.APIKey)
		}
		if userKey == "" {
			userKey = promptSecret("User Key (ETORO_USER_KEY)", cfg.Auth.UserKey)
		}

		if apiKey == "" || userKey == "" {
			return errorf("both API key and user key are required")
		}

		cfg.Auth.APIKey = apiKey
		cfg.Auth.UserKey = userKey

		client := api.NewClient(cfg, flagDemo)
		identity, err := client.GetIdentity()
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{
				"status": "authenticated",
				"gcid":   identity.GCID,
				"config": config.ConfigPath(),
			})
		} else {
			output.Successf("Authenticated successfully (GCID: %d).", identity.GCID)
			output.Infof("Config saved to %s", config.ConfigPath())
		}
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored API keys",
	Long: `Clear your API keys from the config file.

Note: this does not revoke the keys on eToro's servers — it only
removes them from your local config.

Examples:
  etoro auth logout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := config.Load()

		if cfg.Auth.APIKey == "" && cfg.Auth.UserKey == "" {
			if output.GetFormat() == output.JSON {
				output.PrintJSON(map[string]any{"status": "already_logged_out"})
			} else {
				output.Infof("No credentials stored.")
			}
			return nil
		}

		cfg.Auth.APIKey = ""
		cfg.Auth.UserKey = ""

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "logged_out", "config": config.ConfigPath()})
		} else {
			output.Successf("Logged out — credentials removed from %s", config.ConfigPath())
		}
		return nil
	},
}

func promptSecret(label, current string) string {
	if current != "" {
		masked := current
		if len(masked) > 8 {
			masked = masked[:4] + "•••" + masked[len(masked)-4:]
		}
		fmt.Fprintf(os.Stderr, "  %s [%s]: ", label, masked)
	} else {
		fmt.Fprintf(os.Stderr, "  %s: ", label)
	}
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return current
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return current
	}
	return val
}

func init() {
	authLoginCmd.Flags().String("api-key", "", "eToro public API key")
	authLoginCmd.Flags().String("user-key", "", "eToro user key")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
