package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/marianopa-tr/etoro-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var quoteCmd = &cobra.Command{
	Use:   "quote <symbols...>",
	Short: "Get live quotes for one or more instruments",
	Long: `Display current bid/ask prices for instruments. Use --watch for
live auto-refreshing quotes.

Examples:
  etoro quote AAPL
  etoro quote AAPL TSLA GOOG
  etoro quote BTC --watch
  etoro quote AAPL --watch --interval 5s`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		ids, symbols, err := res.ResolveMultiple(args)
		if err != nil {
			return err
		}

		watch, _ := cmd.Flags().GetBool("watch")
		interval, _ := cmd.Flags().GetDuration("interval")

		if watch {
			return watchQuotes(client, ids, symbols, interval)
		}

		return printQuotes(client, ids, symbols)
	},
}

func printQuotes(client *api.Client, ids []int, symbols []string) error {
	instruments, err := client.GetInstruments(ids)
	if err != nil {
		return err
	}

	instMap := make(map[int]api.InstrumentDisplayData)
	for _, inst := range instruments.InstrumentDisplayDatas {
		instMap[inst.InstrumentID] = inst
	}

	rates, err := client.GetRates(ids)
	if err != nil {
		return err
	}

	rows := make([]output.QuoteRow, 0, len(rates.Rates))
	for i, rate := range rates.Rates {
		inst := instMap[rate.InstrumentID]
		symbol := ""
		if i < len(symbols) {
			symbol = symbols[i]
		}
		if inst.Symbol != "" {
			symbol = inst.Symbol
		}
		rows = append(rows, output.QuoteRow{
			Symbol: symbol,
			Name:   inst.DisplayName,
			Bid:    rate.Bid,
			Ask:    rate.Ask,
			Last:   rate.LastExecution,
		})
	}

	output.PrintQuotes(rows, output.GetFormat())
	return nil
}

func watchQuotes(client *api.Client, ids []int, symbols []string, interval time.Duration) error {
	if interval == 0 {
		interval = 2 * time.Second
	}

	instruments, err := client.GetInstruments(ids)
	if err != nil {
		return err
	}

	instMap := make(map[int]api.InstrumentDisplayData)
	for _, inst := range instruments.InstrumentDisplayDatas {
		instMap[inst.InstrumentID] = inst
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	printWatchQuotes(client, ids, symbols, instMap)

	for {
		select {
		case <-sig:
			fmt.Fprintln(os.Stderr, "\nStopped watching.")
			return nil
		case <-ticker.C:
			fmt.Print("\033[H\033[2J")
			printWatchQuotes(client, ids, symbols, instMap)
		}
	}
}

func printWatchQuotes(client *api.Client, ids []int, symbols []string, instMap map[int]api.InstrumentDisplayData) {
	rates, err := client.GetRates(ids)
	if err != nil {
		output.Errorf("failed to fetch rates: %s", err)
		return
	}

	rows := make([]output.QuoteRow, 0, len(rates.Rates))
	for i, rate := range rates.Rates {
		inst := instMap[rate.InstrumentID]
		symbol := ""
		if i < len(symbols) {
			symbol = symbols[i]
		}
		if inst.Symbol != "" {
			symbol = inst.Symbol
		}
		rows = append(rows, output.QuoteRow{
			Symbol: symbol,
			Name:   inst.DisplayName,
			Bid:    rate.Bid,
			Ask:    rate.Ask,
			Last:   rate.LastExecution,
		})
	}

	fmt.Fprintf(os.Stderr, "  Live quotes (every %s) — Ctrl+C to stop\n\n", fmt.Sprint(time.Now().Format("15:04:05")))
	output.PrintQuotes(rows, output.GetFormat())
}

func init() {
	quoteCmd.Flags().Bool("watch", false, "auto-refresh quotes")
	quoteCmd.Flags().Duration("interval", 2*time.Second, "refresh interval (e.g. 2s, 5s)")
	rootCmd.AddCommand(quoteCmd)
}
