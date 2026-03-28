package cmd

import (
	"time"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "View your portfolio, positions, orders, and history",
	Long: `Access your eToro portfolio information including equity summary,
open positions, pending orders, and trade history.

Examples:
  etoro portfolio summary
  etoro portfolio positions
  etoro portfolio orders
  etoro portfolio history --days 30`,
}

var portfolioSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Portfolio summary with P&L and equity",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetPnL()
		if err != nil {
			return err
		}

		p := resp.ClientPortfolio
		summary := output.PortfolioSummary{
			Credit:        p.Credit,
			UnrealizedPnL: p.UnrealizedPnL,
			Equity:        p.Credit + p.UnrealizedPnL,
			PositionCount: len(p.Positions),
			OrderCount:    len(p.Orders) + len(p.OrdersForOpen),
			CopyCount:     len(p.Mirrors),
			IsDemo:        flagDemo,
		}

		output.PrintPortfolioSummary(summary, output.GetFormat())
		return nil
	},
}

var portfolioPositionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "List open positions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetPortfolio()
		if err != nil {
			return err
		}

		positions := resp.ClientPortfolio.Positions
		ids := make([]int, len(positions))
		for i, p := range positions {
			ids[i] = p.InstrumentID
		}
		symbols := resolveInstrumentSymbols(client, ids)

		rows := make([]output.PositionRow, len(positions))
		for i, pos := range positions {
			dir := "Long"
			if !pos.IsBuy {
				dir = "Short"
			}
			rows[i] = output.PositionRow{
				PositionID: pos.PositionID,
				Symbol:     symbols[pos.InstrumentID],
				Direction:  dir,
				Amount:     pos.Amount,
				Units:      pos.Units,
				OpenRate:   pos.OpenRate,
				Leverage:   pos.Leverage,
				SL:         pos.StopLossRate,
				TP:         pos.TakeProfitRate,
				OpenDate:   formatTime(pos.OpenDateTime),
			}
		}

		output.PrintPositions(rows, output.GetFormat())
		return nil
	},
}

var portfolioOrdersCmd = &cobra.Command{
	Use:   "orders",
	Short: "List pending orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetPortfolio()
		if err != nil {
			return err
		}

		p := resp.ClientPortfolio

		var allIDs []int
		for _, o := range p.Orders {
			allIDs = append(allIDs, o.InstrumentID)
		}
		for _, o := range p.OrdersForOpen {
			allIDs = append(allIDs, o.InstrumentID)
		}
		symbols := resolveInstrumentSymbols(client, allIDs)

		rows := make([]output.OrderRow, 0)

		for _, o := range p.Orders {
			dir := "Long"
			if !o.IsBuy {
				dir = "Short"
			}
			rows = append(rows, output.OrderRow{
				OrderID:   o.OrderID,
				Symbol:    symbols[o.InstrumentID],
				Direction: dir,
				Amount:    o.Amount,
				Rate:      o.Rate,
				Leverage:  o.Leverage,
				SL:        o.StopLossRate,
				TP:        o.TakeProfitRate,
				OpenDate:  formatTime(o.OpenDateTime),
			})
		}

		for _, o := range p.OrdersForOpen {
			dir := "Long"
			if !o.IsBuy {
				dir = "Short"
			}
			rows = append(rows, output.OrderRow{
				OrderID:   o.OrderID,
				Symbol:    symbols[o.InstrumentID],
				Direction: dir,
				Amount:    o.Amount,
				Leverage:  o.Leverage,
				SL:        o.StopLossRate,
				TP:        o.TakeProfitRate,
				OpenDate:  formatTime(o.OpenDateTime),
			})
		}

		output.PrintOrders(rows, output.GetFormat())
		return nil
	},
}

var portfolioHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View trade history",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		days, _ := cmd.Flags().GetInt("days")
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		minDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

		entries, err := client.GetTradeHistory(minDate, page, pageSize)
		if err != nil {
			return err
		}

		entryIDs := make([]int, len(entries))
		for i, e := range entries {
			entryIDs[i] = e.InstrumentID
		}
		symbols := resolveInstrumentSymbols(client, entryIDs)

		rows := make([]output.HistoryRow, len(entries))
		for i, e := range entries {
			dir := "Long"
			if !e.IsBuy {
				dir = "Short"
			}
			rows[i] = output.HistoryRow{
				PositionID: e.PositionID,
				Symbol:     symbols[e.InstrumentID],
				Direction:  dir,
				Investment: e.Investment,
				NetProfit:  e.NetProfit,
				OpenRate:   e.OpenRate,
				CloseRate:  e.CloseRate,
				Leverage:   e.Leverage,
				Units:      e.Units,
				Fees:       e.Fees,
				Opened:     formatTime(e.OpenTimestamp),
				Closed:     formatTime(e.CloseTimestamp),
			}
		}

		output.PrintTradeHistory(rows, output.GetFormat())
		return nil
	},
}

func formatTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.Format("2006-01-02 15:04")
}

func init() {
	portfolioHistoryCmd.Flags().Int("days", 30, "number of days of history to show")
	portfolioHistoryCmd.Flags().Int("page", 1, "page number")
	portfolioHistoryCmd.Flags().Int("page-size", 50, "results per page")

	portfolioCmd.AddCommand(portfolioSummaryCmd)
	portfolioCmd.AddCommand(portfolioPositionsCmd)
	portfolioCmd.AddCommand(portfolioOrdersCmd)
	portfolioCmd.AddCommand(portfolioHistoryCmd)
	rootCmd.AddCommand(portfolioCmd)
}
