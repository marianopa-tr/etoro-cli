package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/etoro/etoro-cli/internal/api"
	"github.com/etoro/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var instrumentsCmd = &cobra.Command{
	Use:   "instruments",
	Short: "Get instrument details",
	Long: `Retrieve detailed information about financial instruments.

Examples:
  etoro instruments get AAPL
  etoro instruments get 1001
  etoro instruments get TSLA --output json`,
}

var instrumentsGetCmd = &cobra.Command{
	Use:   "get <symbol|id>",
	Short: "Get detailed information about an instrument",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		query := args[0]

		var inst *api.Instrument

		if id, err := strconv.Atoi(query); err == nil {
			displayResp, err := client.GetInstruments([]int{id})
			if err != nil {
				return err
			}
			if len(displayResp.InstrumentDisplayDatas) == 0 {
				return fmt.Errorf("instrument ID %d not found", id)
			}
			dd := displayResp.InstrumentDisplayDatas[0]
			inst = &api.Instrument{
				InstrumentID: dd.InstrumentID,
				Symbol:       dd.Symbol,
				DisplayName:  dd.DisplayName,
			}
		} else {
			resp, err := client.SearchInstruments(query, 1, 10)
			if err != nil {
				return err
			}
			for i, item := range resp.Items {
				if strings.EqualFold(item.Symbol, query) {
					inst = &resp.Items[i]
					break
				}
			}
			if inst == nil && len(resp.Items) > 0 {
				inst = &resp.Items[0]
			}
			if inst == nil {
				return fmt.Errorf("instrument %q not found", query)
			}
		}

		detail := output.InstrumentDetail{
			ID:            inst.InstrumentID,
			Symbol:        inst.Symbol,
			Name:          inst.DisplayName,
			Type:          inst.InternalAssetClassName,
			Exchange:      inst.InternalExchangeName,
			AssetClass:    inst.InternalAssetClassName,
			Price:         inst.CurrentRate,
			DailyChg:      inst.DailyPriceChange,
			WeeklyChg:     inst.WeeklyPriceChange,
			MonthlyChg:    inst.MonthlyPriceChange,
			ThreeMonthChg: inst.ThreeMonthPriceChange,
			OneYearChg:    inst.OneYearPriceChange,
			Open:          inst.IsCurrentlyTradable,
			Tradable:      inst.IsCurrentlyTradable,
			BuyEnabled:    inst.IsBuyEnabled,
			Popularity:    inst.Popularity7Day,
		}

		output.PrintInstrumentDetail(detail, output.GetFormat())
		return nil
	},
}

func init() {
	instrumentsCmd.AddCommand(instrumentsGetCmd)
	rootCmd.AddCommand(instrumentsCmd)
}
