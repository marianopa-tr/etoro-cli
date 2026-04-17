# etoro-cli

Trade, invest, and copy from your terminal. Beautiful colored tables by default, `--output json` everywhere for scripts and AI agents.

## Install

```bash
npm install -g etoro-cli
```

## Usage

```bash
etoro setup                          # configure API keys
etoro search apple                   # search instruments
etoro quote AAPL TSLA                # live quotes
etoro portfolio summary              # portfolio overview
etoro trade open AAPL --amount 500   # open a position
etoro shell                          # interactive REPL
```

Every command supports `--output json` for machine-readable output.

## Platforms

| Package | OS | Architecture |
|---|---|---|
| etoro-cli-darwin-arm64 | macOS | Apple Silicon |
| etoro-cli-darwin-x64 | macOS | Intel |
| etoro-cli-linux-arm64 | Linux | ARM64 |
| etoro-cli-linux-x64 | Linux | x86_64 |
| etoro-cli-windows-arm64 | Windows | ARM64 |
| etoro-cli-windows-x64 | Windows | x86_64 |

The correct binary is selected automatically at install time via `optionalDependencies`.

## Documentation

Full documentation, configuration options, and examples: [github.com/marianopa-tr/etoro-cli](https://github.com/marianopa-tr/etoro-cli)

## License

MIT
