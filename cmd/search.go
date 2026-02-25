package cmd

import (
	"strings"

	"github.com/etoro/etoro-cli/internal/api"
	"github.com/etoro/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for instruments by name or symbol",
	Long: `Search the eToro instrument catalog. Results include current price,
daily performance, and trading status.

Examples:
  etoro search apple
  etoro search BTC
  etoro search "S&P 500" --page-size 5`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		query := strings.Join(args, " ")

		resp, err := client.SearchInstrumentsByText(query, page, pageSize)
		if err != nil {
			return err
		}

		rows := make([]output.InstrumentRow, 0, len(resp.Items))
		for _, inst := range resp.Items {
			if inst.IsHiddenFromClient || inst.InstrumentID < 0 {
				continue
			}
			rows = append(rows, instrumentToRow(inst))
		}

		output.PrintInstrumentSearch(rows, resp.TotalItems, output.GetFormat())
		return nil
	},
}

func instrumentToRow(inst api.Instrument) output.InstrumentRow {
	return output.InstrumentRow{
		ID:         inst.InstrumentID,
		Symbol:     inst.Symbol,
		Name:       inst.DisplayName,
		Type:       inst.InternalAssetClassName,
		Exchange:   inst.InternalExchangeName,
		Price:      inst.CurrentRate,
		DailyChg:   inst.DailyPriceChange,
		WeeklyChg:  inst.WeeklyPriceChange,
		Open:       inst.IsExchangeOpen,
		Tradable:   inst.IsCurrentlyTradable,
		Popularity: inst.Popularity7Day,
	}
}

func init() {
	searchCmd.Flags().Int("page", 1, "page number")
	searchCmd.Flags().Int("page-size", 20, "results per page")
	rootCmd.AddCommand(searchCmd)
}
