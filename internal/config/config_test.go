package config

import (
	"os"
	"path/filepath"
	"testing"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

func TestLoadDefaults(t *testing.T) {
	chdir(t, t.TempDir())
	t.Setenv("ETORO_API_KEY", "")
	t.Setenv("ETORO_PUBLIC_KEY", "")
	t.Setenv("ETORO_USER_KEY", "")
	t.Setenv("ETORO_CONFIG_DIR", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Defaults.Output != "table" {
		t.Errorf("Output = %q, want %q", cfg.Defaults.Output, "table")
	}
	if cfg.Defaults.Timeout != "30s" {
		t.Errorf("Timeout = %q, want %q", cfg.Defaults.Timeout, "30s")
	}
	if cfg.Auth.APIKey != "" {
		t.Errorf("APIKey = %q, want empty", cfg.Auth.APIKey)
	}
}

func TestLoadFromEnvVars(t *testing.T) {
	chdir(t, t.TempDir())
	t.Setenv("ETORO_CONFIG_DIR", t.TempDir())
	t.Setenv("ETORO_PUBLIC_KEY", "")
	t.Setenv("ETORO_API_KEY", "test-api-key")
	t.Setenv("ETORO_USER_KEY", "test-user-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Auth.APIKey != "test-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.Auth.APIKey, "test-api-key")
	}
	if cfg.Auth.UserKey != "test-user-key" {
		t.Errorf("UserKey = %q, want %q", cfg.Auth.UserKey, "test-user-key")
	}
}

func TestLoadFromTOMLFile(t *testing.T) {
	chdir(t, t.TempDir())
	dir := t.TempDir()
	t.Setenv("ETORO_CONFIG_DIR", dir)
	t.Setenv("ETORO_API_KEY", "")
	t.Setenv("ETORO_PUBLIC_KEY", "")
	t.Setenv("ETORO_USER_KEY", "")

	tomlContent := `
[auth]
api_key = "toml-api-key"
user_key = "toml-user-key"

[defaults]
output = "json"
timeout = "10s"
demo = true
`
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(tomlContent), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Auth.APIKey != "toml-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.Auth.APIKey, "toml-api-key")
	}
	if cfg.Auth.UserKey != "toml-user-key" {
		t.Errorf("UserKey = %q, want %q", cfg.Auth.UserKey, "toml-user-key")
	}
	if cfg.Defaults.Output != "json" {
		t.Errorf("Output = %q, want %q", cfg.Defaults.Output, "json")
	}
	if cfg.Defaults.Timeout != "10s" {
		t.Errorf("Timeout = %q, want %q", cfg.Defaults.Timeout, "10s")
	}
	if !cfg.Defaults.Demo {
		t.Error("Demo = false, want true")
	}
}

func TestEnvOverridesToml(t *testing.T) {
	chdir(t, t.TempDir())
	dir := t.TempDir()
	t.Setenv("ETORO_CONFIG_DIR", dir)
	t.Setenv("ETORO_PUBLIC_KEY", "")

	tomlContent := `
[auth]
api_key = "toml-key"
user_key = "toml-user"
`
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(tomlContent), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ETORO_API_KEY", "env-key")
	t.Setenv("ETORO_USER_KEY", "env-user")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Auth.APIKey != "env-key" {
		t.Errorf("APIKey = %q, want %q (env should override TOML)", cfg.Auth.APIKey, "env-key")
	}
	if cfg.Auth.UserKey != "env-user" {
		t.Errorf("UserKey = %q, want %q (env should override TOML)", cfg.Auth.UserKey, "env-user")
	}
}

func TestLoadEnvFile(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := `# comment
TEST_LOAD_ENV_FILE_KEY=hello-from-env-file
TEST_LOAD_ENV_FILE_QUOTED="quoted-value"
`
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	os.Unsetenv("TEST_LOAD_ENV_FILE_KEY")
	os.Unsetenv("TEST_LOAD_ENV_FILE_QUOTED")

	loadEnvFile(envPath)

	if v := os.Getenv("TEST_LOAD_ENV_FILE_KEY"); v != "hello-from-env-file" {
		t.Errorf("env key = %q, want %q", v, "hello-from-env-file")
	}
	if v := os.Getenv("TEST_LOAD_ENV_FILE_QUOTED"); v != "quoted-value" {
		t.Errorf("env quoted = %q, want %q", v, "quoted-value")
	}

	os.Unsetenv("TEST_LOAD_ENV_FILE_KEY")
	os.Unsetenv("TEST_LOAD_ENV_FILE_QUOTED")
}

func TestLoadEnvFileDoesNotOverrideExisting(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := `TEST_NO_OVERRIDE=from-file`
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TEST_NO_OVERRIDE", "from-env")

	loadEnvFile(envPath)

	if v := os.Getenv("TEST_NO_OVERRIDE"); v != "from-env" {
		t.Errorf("env = %q, want %q (existing should not be overridden)", v, "from-env")
	}
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		userKey string
		wantErr bool
	}{
		{"both set", "key", "user", false},
		{"missing api key", "", "user", true},
		{"missing user key", "key", "", true},
		{"both missing", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Auth: AuthConfig{APIKey: tt.apiKey, UserKey: tt.userKey},
			}
			err := cfg.RequireAuth()
			if (err != nil) != tt.wantErr {
				t.Errorf("RequireAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeoutDuration(t *testing.T) {
	tests := []struct {
		name    string
		timeout string
		wantSec int
	}{
		{"default", "", 30},
		{"10 seconds", "10s", 10},
		{"1 minute", "1m", 60},
		{"invalid fallback", "invalid", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Defaults: DefaultsConfig{Timeout: tt.timeout}}
			d := cfg.TimeoutDuration()
			got := int(d.Seconds())
			if got != tt.wantSec {
				t.Errorf("TimeoutDuration() = %ds, want %ds", got, tt.wantSec)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	chdir(t, t.TempDir())
	dir := t.TempDir()
	t.Setenv("ETORO_CONFIG_DIR", dir)
	t.Setenv("ETORO_API_KEY", "")
	t.Setenv("ETORO_PUBLIC_KEY", "")
	t.Setenv("ETORO_USER_KEY", "")

	cfg := &Config{
		Auth: AuthConfig{APIKey: "saved-key", UserKey: "saved-user"},
		Defaults: DefaultsConfig{
			Output:  "json",
			Timeout: "5s",
			Demo:    true,
		},
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Auth.APIKey != "saved-key" {
		t.Errorf("APIKey = %q, want %q", loaded.Auth.APIKey, "saved-key")
	}
	if loaded.Auth.UserKey != "saved-user" {
		t.Errorf("UserKey = %q, want %q", loaded.Auth.UserKey, "saved-user")
	}
	if loaded.Defaults.Output != "json" {
		t.Errorf("Output = %q, want %q", loaded.Defaults.Output, "json")
	}
}

func TestConfigDir(t *testing.T) {
	t.Setenv("ETORO_CONFIG_DIR", "/custom/path")
	if got := ConfigDir(); got != "/custom/path" {
		t.Errorf("ConfigDir() = %q, want %q", got, "/custom/path")
	}
}

func TestConfigDirDefault(t *testing.T) {
	t.Setenv("ETORO_CONFIG_DIR", "")
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir() returned empty string")
	}
	if filepath.Base(dir) != "etoro" {
		t.Errorf("ConfigDir() = %q, want basename 'etoro'", dir)
	}
}
