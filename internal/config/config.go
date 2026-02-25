package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Auth     AuthConfig     `toml:"auth"`
	Defaults DefaultsConfig `toml:"defaults"`
}

type AuthConfig struct {
	APIKey  string `toml:"api_key"`
	UserKey string `toml:"user_key"`
}

type DefaultsConfig struct {
	Output  string `toml:"output"`
	Demo    bool   `toml:"demo"`
	Timeout string `toml:"timeout"`
}

const (
	DefaultBaseURL = "https://public-api.etoro.com"
	DefaultTimeout = 30 * time.Second
)

func Load() (*Config, error) {
	cfg := &Config{
		Defaults: DefaultsConfig{
			Output:  "table",
			Timeout: "30s",
		},
	}

	path := ConfigPath()
	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, cfg); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}
	}

	loadEnvFile(".env")
	if p := os.Getenv("ETORO_ENV_PATH"); p != "" {
		loadEnvFile(p)
	}

	if v := os.Getenv("ETORO_API_KEY"); v != "" {
		cfg.Auth.APIKey = v
	}
	if v := os.Getenv("ETORO_PUBLIC_KEY"); v != "" {
		cfg.Auth.APIKey = v
	}
	if v := os.Getenv("ETORO_USER_KEY"); v != "" {
		cfg.Auth.UserKey = v
	}

	return cfg, nil
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func (c *Config) RequireAuth() error {
	if c.Auth.APIKey == "" {
		return fmt.Errorf("API key not configured. Run `etoro setup` or set ETORO_PUBLIC_KEY")
	}
	if c.Auth.UserKey == "" {
		return fmt.Errorf("user key not configured. Run `etoro setup` or set ETORO_USER_KEY")
	}
	return nil
}

func (c *Config) TimeoutDuration() time.Duration {
	if c.Defaults.Timeout == "" {
		return DefaultTimeout
	}
	d, err := time.ParseDuration(c.Defaults.Timeout)
	if err != nil {
		return DefaultTimeout
	}
	return d
}

func ConfigDir() string {
	if v := os.Getenv("ETORO_CONFIG_DIR"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "etoro")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.toml")
}

func CacheDir() string {
	if v := os.Getenv("ETORO_CACHE_DIR"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "etoro")
}

func Save(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(ConfigPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
