package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func skipIfNoKeys(t *testing.T) {
	t.Helper()
	apiKey := os.Getenv("ETORO_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ETORO_PUBLIC_KEY")
	}
	if apiKey == "" || os.Getenv("ETORO_USER_KEY") == "" {
		t.Skip("ETORO_PUBLIC_KEY and ETORO_USER_KEY required for live E2E tests")
	}
	if os.Getenv("ETORO_LIVE_TESTS") != "1" {
		t.Skip("Set ETORO_LIVE_TESTS=1 to run live API tests")
	}
}

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := t.TempDir() + "/etoro"
	cmd := exec.Command("go", "build", "-o", binary, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %s\n%s", err, out)
	}
	return binary
}

func runCLI(t *testing.T, binary string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("exec error: %v", err)
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

func TestE2E_Help(t *testing.T) {
	bin := buildCLI(t)

	stdout, _, code := runCLI(t, bin, "--help")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "etoro") {
		t.Errorf("help should mention 'etoro'")
	}
}

func TestE2E_HelpSubcommands(t *testing.T) {
	bin := buildCLI(t)

	subcommands := []string{
		"search", "quote", "instruments", "portfolio", "trade",
		"orders", "watchlist", "copy", "feed", "pi", "auth",
		"status", "setup", "shell", "upgrade",
	}

	for _, sub := range subcommands {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("help output is empty")
			}
		})
	}
}

func TestE2E_VersionFlags(t *testing.T) {
	bin := buildCLI(t)

	tests := []struct {
		args []string
	}{
		{[]string{"--output", "json", "--help"}},
		{[]string{"--demo", "--help"}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, "_"), func(t *testing.T) {
			_, _, code := runCLI(t, bin, tt.args...)
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
		})
	}
}

func TestE2E_StatusNoAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "status")
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	if !strings.Contains(stdout.String(), "missing") {
		t.Logf("output: %s", stdout.String())
	}
}

func TestE2E_AuthStatusNoAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "auth", "status")
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	if !strings.Contains(stdout.String(), "missing") {
		t.Logf("output: %s", stdout.String())
	}
}

func TestE2E_StatusJSON(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "status", "--output", "json")
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
	out := stdout.String()
	if !strings.Contains(out, "apiKey") {
		t.Errorf("JSON output should contain 'apiKey', got %q", out)
	}
}

func TestE2E_SearchRequiresAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "search", "AAPL")
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Error("expected error when no auth configured")
	}
	if !strings.Contains(stderr.String(), "API key not configured") {
		t.Logf("stderr: %s", stderr.String())
	}
}

func TestE2E_QuoteRequiresAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "quote", "AAPL")
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Error("expected error when no auth configured")
	}
}

func TestE2E_TradeRequiresAmount(t *testing.T) {
	bin := buildCLI(t)
	_, stderr, code := runCLI(t, bin, "trade", "open", "AAPL")
	if code == 0 {
		t.Error("expected error when no --amount or --units")
	}
	_ = stderr
}

func TestE2E_SearchMultipleWords(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "search", "S&P", "500", "--help")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Run()
}

func TestE2E_PortfolioSubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"summary", "positions", "orders", "history"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "portfolio", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_TradeSubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"open", "close", "limit"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "trade", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_WatchlistSubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"list", "create", "add", "remove"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "watchlist", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_CopySubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"discover", "performance", "copiers"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "copy", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_FeedSubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"list", "post"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "feed", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_PISubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"copiers", "get"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "pi", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_OrdersSubcommands(t *testing.T) {
	bin := buildCLI(t)
	subs := []string{"list", "cancel", "cancel-all"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			stdout, _, code := runCLI(t, bin, "orders", sub, "--help")
			if code != 0 {
				t.Errorf("exit code = %d", code)
			}
			if stdout == "" {
				t.Error("empty help")
			}
		})
	}
}

func TestE2E_InvalidCommand(t *testing.T) {
	bin := buildCLI(t)
	_, stderr, code := runCLI(t, bin, "nonexistent")
	if code == 0 {
		t.Error("expected non-zero exit code for invalid command")
	}
	_ = stderr
}

func TestE2E_WatchlistAddRequiresAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "watchlist", "add", "AAPL")
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Error("expected error when no auth configured")
	}
}

func TestE2E_WatchlistRemoveRequiresAuth(t *testing.T) {
	bin := buildCLI(t)
	cmd := exec.Command(bin, "watchlist", "remove", "AAPL")
	cmd.Env = []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Error("expected error when no auth configured")
	}
}

func TestE2E_TradeLimitRequiresPrice(t *testing.T) {
	bin := buildCLI(t)
	_, stderr, code := runCLI(t, bin, "trade", "limit", "AAPL", "--amount", "500")
	if code == 0 {
		t.Error("expected error when --price missing")
	}
	_ = stderr
}

func TestE2E_FeedListRequiresFilter(t *testing.T) {
	bin := buildCLI(t)
	_, stderr, code := runCLI(t, bin, "feed", "list")
	if code == 0 {
		t.Error("expected error when no --instrument or --username")
	}
	_ = stderr
}

// --- Live API tests (require ETORO_API_KEY and ETORO_USER_KEY) ---

func TestE2E_Live_Search(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "search", "Apple", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "AAPL") && !strings.Contains(stdout, "Apple") {
		t.Errorf("search output should contain AAPL or Apple, got: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_SearchMultiWord(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "search", "Apple", "Inc", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "Apple") {
		t.Logf("output: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_Quote(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "quote", "AAPL", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "bid") && !strings.Contains(stdout, "Bid") {
		t.Logf("quote output: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_QuoteMultiple(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "quote", "AAPL", "TSLA", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_InstrumentsGet(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "instruments", "get", "AAPL", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "AAPL") && !strings.Contains(stdout, "Apple") {
		t.Logf("output: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_InstrumentsGetByID(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "instruments", "get", "1001", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_Status(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "status", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "connected") {
		t.Logf("output: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_AuthStatus(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "auth", "status", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	if !strings.Contains(stdout, "authenticated") {
		t.Logf("output: %s", stdout[:min(200, len(stdout))])
	}
}

func TestE2E_Live_PortfolioSummary(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "portfolio", "summary", "--demo", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_PortfolioPositions(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "portfolio", "positions", "--demo", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_PortfolioOrders(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "portfolio", "orders", "--demo", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_PortfolioHistory(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "portfolio", "history", "--demo", "--days", "30", "--output", "json")
	if code != 0 && !strings.Contains(stdout, "InsufficientPermissions") && !strings.Contains(stderr, "InsufficientPermissions") {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_OrdersList(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "orders", "list", "--demo", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_WatchlistList(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "watchlist", "list", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_WatchlistCurated(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "watchlist", "list", "--curated", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_CopyDiscover(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "copy", "discover", "--period", "LastYear", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_PIcopiers(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "pi", "copiers", "--output", "json")
	if code != 0 && !strings.Contains(stdout, "InsufficientPermissions") && !strings.Contains(stderr, "InsufficientPermissions") {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_FeedInstrument(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "feed", "list", "--instrument", "AAPL", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_SearchPageSize(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "search", "tech", "--page-size", "3", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func TestE2E_Live_DemoPortfolio(t *testing.T) {
	skipIfNoKeys(t)
	bin := buildCLI(t)
	stdout, stderr, code := runCLI(t, bin, "portfolio", "summary", "--demo", "--output", "json")
	if code != 0 {
		t.Errorf("exit %d: %s", code, stderr)
	}
	_ = stdout
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
