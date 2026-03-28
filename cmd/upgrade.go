package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

const (
	repoOwner = "etoro"
	repoName  = "etoro-cli"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Check for updates and self-update the CLI",
	Long: `Check for newer versions of the eToro CLI and optionally upgrade
to the latest release from GitHub.

Examples:
  etoro upgrade
  etoro upgrade --check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")

		output.Infof("Current version: %s", version)
		output.Infof("Checking for updates...")

		latest, downloadURL, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("checking for updates: %w", err)
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{
				"current":    version,
				"latest":     latest,
				"updateAvailable": latest != version,
				"downloadUrl":     downloadURL,
			})
			return nil
		}

		if latest == version {
			output.Successf("You're on the latest version (%s).", version)
			return nil
		}

		output.Infof("New version available: %s → %s", output.Yellow(version), output.Green(latest))

		if checkOnly {
			output.Infof("Run `etoro upgrade` to update.")
			return nil
		}

		if !flagYes {
			if !output.Confirm(fmt.Sprintf("Upgrade to %s?", latest)) {
				output.Infof("Upgrade cancelled.")
				return nil
			}
		}

		if downloadURL == "" {
			output.Infof("Download the latest release from:")
			output.Infof("  https://github.com/%s/%s/releases/latest", repoOwner, repoName)
			return nil
		}

		output.Infof("Downloading %s...", latest)
		if err := downloadAndReplace(downloadURL); err != nil {
			return fmt.Errorf("upgrading: %w", err)
		}

		output.Successf("Upgraded to %s!", latest)
		return nil
	},
}

func getLatestRelease() (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	osName := runtime.GOOS
	archName := runtime.GOARCH

	var downloadURL string
	for _, asset := range release.Assets {
		name := asset.Name
		if matchesPlatform(name, osName, archName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	return release.TagName, downloadURL, nil
}

func matchesPlatform(name, osName, archName string) bool {
	osMatches := false
	archMatches := false

	switch osName {
	case "darwin":
		osMatches = contains(name, "darwin") || contains(name, "macos")
	case "linux":
		osMatches = contains(name, "linux")
	case "windows":
		osMatches = contains(name, "windows")
	}

	switch archName {
	case "amd64":
		archMatches = contains(name, "amd64") || contains(name, "x86_64")
	case "arm64":
		archMatches = contains(name, "arm64") || contains(name, "aarch64")
	}

	return osMatches && archMatches
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && findSubstring(s, substr))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func downloadAndReplace(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable path: %w", err)
	}

	tmpFile := execPath + ".new"
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpFile)
		return fmt.Errorf("downloading: %w", err)
	}
	f.Close()

	oldFile := execPath + ".old"
	if err := os.Rename(execPath, oldFile); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("backing up current binary: %w", err)
	}

	if err := os.Rename(tmpFile, execPath); err != nil {
		os.Rename(oldFile, execPath) // rollback
		return fmt.Errorf("replacing binary: %w", err)
	}

	os.Remove(oldFile)
	return nil
}

func init() {
	upgradeCmd.Flags().Bool("check", false, "only check for updates, don't install")
	rootCmd.AddCommand(upgradeCmd)
}
