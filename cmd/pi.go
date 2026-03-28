package cmd

import (
	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var piCmd = &cobra.Command{
	Use:   "pi",
	Short: "Popular Investor program data",
	Long: `Access Popular Investor (PI) data including copier information
and trader analytics.

Examples:
  etoro pi copiers
  etoro pi get BigTech
  etoro pi get trader123 --period CurrYear
  etoro pi gain trader123`,
}

var piCopiersCmd = &cobra.Command{
	Use:   "copiers",
	Short: "View your copier details (PI data)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetCopiers()
		if err != nil {
			return err
		}

		rows := make([]output.CopierRow, len(resp.Copiers))
		for i, c := range resp.Copiers {
			rows[i] = output.CopierRow{
				Country:  c.Country,
				Club:     c.Club,
				Duration: c.CopyStartedCategory,
				Amount:   c.AmountCategory,
			}
		}

		output.PrintCopiers(rows, output.GetFormat())
		return nil
	},
}

var piGetCmd = &cobra.Command{
	Use:   "get <username>",
	Short: "Get trader profile and trading stats",
	Long: `Get detailed trader profile including performance, risk metrics, and statistics.

Periods: CurrMonth, CurrQuarter, CurrYear, LastYear, LastTwoYears,
         OneMonthAgo, TwoMonthsAgo, ThreeMonthsAgo, SixMonthsAgo, OneYearAgo

Examples:
  etoro pi get BigTech
  etoro pi get trader123 --period CurrYear`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		period, _ := cmd.Flags().GetString("period")

		info, err := client.GetUserTradeInfo(args[0], period)
		if err != nil {
			return err
		}

		profile := output.TraderProfileRow{
			Username:          info.UserName,
			FullName:          info.FullName,
			IsPopularInvestor: info.IsPopularInvestor,
			Gain:              info.Gain,
			DailyGain:         info.DailyGain,
			WeeklyGain:        info.ThisWeekGain,
			RiskScore:         info.RiskScore,
			MaxDailyRisk:      info.MaxDailyRiskScore,
			MaxMonthlyRisk:    info.MaxMonthlyRiskScore,
			Copiers:           info.Copiers,
			Trades:            info.Trades,
			WinRatio:          info.WinRatio,
			DailyDD:           info.DailyDD,
			WeeklyDD:          info.WeeklyDD,
			PeakToValley:      info.PeakToValley,
			ProfitWeeksPct:    info.ProfitableWeeksPct,
			ProfitMonthsPct:   info.ProfitableMonthsPct,
			AvgPosSize:        info.AvgPosSize,
			AUMTierDesc:       info.AUMTierDesc,
			WeeksRegistered:   info.WeeksSinceRegistration,
		}

		output.PrintTraderProfile(profile, output.GetFormat())
		return nil
	},
}

var piGainCmd = &cobra.Command{
	Use:   "gain <username>",
	Short: "View monthly and yearly gain history for a trader",
	Long: `Display historical monthly and yearly performance data for a trader.

Examples:
  etoro pi gain BigTech
  etoro pi gain trader123 --output json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		gainResp, err := client.GetUserGain(args[0])
		if err != nil {
			return err
		}

		monthly := make([]output.GainRow, len(gainResp.Monthly))
		for i, g := range gainResp.Monthly {
			monthly[i] = output.GainRow{Period: g.Timestamp, Gain: g.Gain}
		}

		yearly := make([]output.GainRow, len(gainResp.Yearly))
		for i, g := range gainResp.Yearly {
			yearly[i] = output.GainRow{Period: g.Timestamp, Gain: g.Gain}
		}

		output.PrintUserGain(args[0], monthly, yearly, output.GetFormat())
		return nil
	},
}

func init() {
	piGetCmd.Flags().String("period", "LastTwoYears", "time period (CurrMonth, CurrQuarter, CurrYear, LastYear, LastTwoYears)")
	piCmd.AddCommand(piCopiersCmd)
	piCmd.AddCommand(piGetCmd)
	piCmd.AddCommand(piGainCmd)
	rootCmd.AddCommand(piCmd)
}
