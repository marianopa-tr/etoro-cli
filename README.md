# eToro CLI

Trade, invest, and copy from your terminal. Beautiful colored tables by default, `--output json` everywhere for scripts and AI agents.

[![CI](https://github.com/marianopa-tr/etoro-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/marianopa-tr/etoro-cli/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/etoro-cli)](https://www.npmjs.com/package/etoro-cli)

## Quick Start

### Install via npm (recommended)

```bash
npm install -g etoro-cli
```

### One-liner install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/marianopa-tr/etoro-cli/main/scripts/install.sh | bash
```

### Install from source

```bash
go install github.com/marianopa-tr/etoro-cli@latest
```

### Build locally

```bash
git clone https://github.com/marianopa-tr/etoro-cli.git
cd etoro-cli
make build      # builds ./etoro with version info
make install    # copies to ~/.local/bin
```

### Get started

```bash
etoro setup        # configure API keys
etoro status       # check connectivity
etoro version      # check installed version
```

## Configuration

Config lives in `~/.config/etoro/config.toml`:

```toml
[auth]
api_key = "your-api-key"
user_key = "your-user-key"

[defaults]
output = "table"     # or "json"
demo = false
timeout = "30s"
```

Resolution order: **CLI flags** > **environment variables** > **config file**.

| Environment Variable    | Description         |
|------------------------|---------------------|
| `ETORO_API_KEY`        | API key             |
| `ETORO_USER_KEY`       | User key            |
| `ETORO_CONFIG_DIR`     | Config directory    |

## Commands

### Market Data

```bash
# Search instruments
etoro search apple
etoro search BTC --page-size 5

# Get instrument details
etoro instruments get AAPL
etoro instruments get 1001

# Live quotes (auto-refreshing)
etoro quote AAPL TSLA GOOG
etoro quote BTC --watch --interval 5s
```

### Portfolio

```bash
# Overview with P&L
etoro portfolio summary

# Open positions
etoro portfolio positions

# Pending orders
etoro portfolio orders

# Trade history (last 30 days)
etoro portfolio history --days 30
```

### Trading

```bash
# Open a position (by amount or units)
etoro trade open AAPL --amount 500 --leverage 2
etoro trade open BTC --units 0.1 --demo

# With stop-loss and take-profit
etoro trade open TSLA --amount 1000 --sl 180 --tp 250

# Short sell
etoro trade open EURUSD --amount 500 --short

# Close a position
etoro trade close 12345

# Place a limit order
etoro trade limit AAPL --price 170 --amount 500
```

### Orders

```bash
etoro orders list
etoro orders cancel 12345
etoro orders cancel-all
```

### Social & Copy Trading

```bash
# Discover top traders
etoro copy discover --period LastYear

# View trader performance
etoro copy performance trader123

# View your copiers
etoro copy copiers

# Social feed
etoro feed list --instrument AAPL
etoro feed post "Bullish on tech this quarter!"
```

### Watchlists

```bash
etoro watchlist list
etoro watchlist list --curated
etoro watchlist create "My Tech Stocks"
etoro watchlist add AAPL --to <watchlistId>
etoro watchlist remove TSLA --from <watchlistId>
```

### Popular Investor

```bash
etoro pi copiers
etoro pi get trader123
```

## Global Flags

| Flag              | Description                           |
|-------------------|---------------------------------------|
| `--demo`          | Use demo/virtual account              |
| `-o, --output`    | Output format: `table` or `json`      |
| `--yes`           | Skip confirmation prompts             |
| `--timeout`       | Request timeout (e.g. `30s`, `1m`)    |

## Agent & Script Usage

Every command supports `--output json` for machine-readable output:

```bash
# Get portfolio as JSON
etoro portfolio summary --output json

# Search and pipe to jq
etoro search AAPL --output json | jq '.instruments[0].instrumentId'

# Automated trading (skip confirmations)
etoro trade open AAPL --amount 500 --demo --yes --output json
```

## Interactive Shell

Start a REPL with history and tab-completion:

```bash
etoro shell
```

```

   ╔██████╗   ╔████████╗   ╔██████╗   ╔██████╗   ╔██████╗
   ██╔═══██╗  ╚═══██╔══╝  ██╔═══██╗  ██╔═══██╗  ██╔═══██╗
   ████████║      ██║     ██║   ██║  ██║   ╚═╝  ██║   ██║
   ██╔═════╝      ██║     ██║   ██║  ██║        ██║   ██║
   ╚███████╗      ██║     ╚██████╔╝  ██║        ╚██████╔╝
    ╚══════╝      ╚═╝      ╚═════╝   ╚═╝         ╚═════╝

                    Interactive Shell

  Type 'help' for commands, 'exit' to quit.

etoro> search apple
etoro> quote AAPL
etoro> portfolio summary
```

## Self-Update

```bash
# Check for updates
etoro upgrade --check

# Upgrade to latest
etoro upgrade
```

## Shell Completions

```bash
# Bash
etoro completion bash > /etc/bash_completion.d/etoro

# Zsh
etoro completion zsh > "${fpath[1]}/_etoro"

# Fish
etoro completion fish > ~/.config/fish/completions/etoro.fish
```

## Development

```bash
make build      # compile with version injection
make test       # run all tests
make vet        # static analysis
make clean      # remove build artifacts
make release    # cross-compile tarballs + checksums into dist/
```

## Releasing

Releases are fully automated. Push a version tag and CI handles everything:

```bash
git tag v0.2.0
git push origin v0.2.0
```

This will:
1. Run tests
2. Cross-compile Go binaries via GoReleaser
3. Create a GitHub Release with assets and checksums
4. Publish all npm packages with provenance attestation

**One-time setup:** add an `NPM_TOKEN` secret to the repository (Settings > Secrets and variables > Actions).

## Project Structure

```
cmd/                  # Cobra command definitions (one file per command group)
internal/
  api/                # Typed HTTP client for eToro API
  config/             # TOML config + env var resolution
  output/             # Rendering layer (tables + JSON)
  resolver/           # Symbol → instrumentId resolution with cache
  shell/              # Interactive REPL
npm/
  cli/                # Meta package (etoro-cli) with platform resolution
  cli-darwin-arm64/   # macOS Apple Silicon binary
  cli-darwin-x64/     # macOS Intel binary
  cli-linux-arm64/    # Linux ARM64 binary
  cli-linux-x64/      # Linux x64 binary
  cli-windows-arm64/  # Windows ARM64 binary
  cli-windows-x64/    # Windows x64 binary
  publish.sh          # Manual npm publish script
  bump-version.sh     # Version bump across all package.json files
scripts/
  install.sh          # Customer-facing binary installer
Makefile              # Build, test, release targets
```

## License

See [LICENSE](LICENSE) for details.
