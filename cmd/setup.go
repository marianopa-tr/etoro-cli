package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/marianopa-tr/etoro-cli/internal/config"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for eToro CLI",
	Long: `Configure your eToro CLI with an interactive, step-by-step wizard.
This will guide you through setting up your API key, user key, and
default preferences.

Examples:
  etoro setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output.PrintBanner("")

		output.PrintNoticeBox("eToro Public API — configure your CLI access.")

		cfg, _ := config.Load()

		// ── Step 1: API Key ─────────────────────────────

		output.PrintStepHeader(1, 3, "API Key")

		dim := color.New(color.Faint)
		dim.Fprintln(os.Stderr, "  Get your API key from the eToro developer portal.")
		dim.Fprintln(os.Stderr, "  https://www.etoro.com/settings/trade")
		fmt.Fprintln(os.Stderr)

		apiKey := wizardPrompt("API Key", cfg.Auth.APIKey, true)
		cfg.Auth.APIKey = apiKey

		// ── Step 2: User Key ────────────────────────────

		output.PrintStepHeader(2, 3, "User Key")

		dim.Fprintln(os.Stderr, "  Your user-specific authentication key.")
		fmt.Fprintln(os.Stderr)

		userKey := wizardPrompt("User Key", cfg.Auth.UserKey, true)
		cfg.Auth.UserKey = userKey

		// ── Step 3: Preferences ─────────────────────────

		output.PrintStepHeader(3, 3, "Preferences")

		outputFmt := wizardPrompt("Default output (table/json)", cfg.Defaults.Output, false)
		cfg.Defaults.Output = outputFmt

		// ── Save ────────────────────────────────────────

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		output.PrintSavedBox(config.ConfigPath())

		dim.Fprintln(os.Stderr, "  Next: run "+output.Cyan("etoro status")+" to verify your connection.")
		fmt.Fprintln(os.Stderr)

		return nil
	},
}

func wizardPrompt(label, defaultVal string, isSecret bool) string {
	green := color.New(color.FgGreen, color.Bold)
	dim := color.New(color.Faint)

	green.Fprint(os.Stderr, "  ❯ ")
	fmt.Fprint(os.Stderr, label)

	if defaultVal != "" {
		displayed := defaultVal
		if isSecret {
			displayed = maskSecret(defaultVal)
		}
		dim.Fprintf(os.Stderr, " [%s]", displayed)
	}

	fmt.Fprint(os.Stderr, ": ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func maskSecret(s string) string {
	if len(s) <= 8 {
		return "••••"
	}
	return s[:4] + "•••" + s[len(s)-4:]
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
