package cmd

import (
	"encoding/json"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Discover and analyze traders to copy",
	Long: `Explore eToro's social trading features: discover top traders,
analyze their performance, and view your copiers.

Examples:
  etoro copy discover --period LastYear
  etoro copy performance trader123
  etoro copy copiers`,
}

var copyDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover traders to copy",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		period, _ := cmd.Flags().GetString("period")
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		data, err := client.SearchUsers(period, page, pageSize)
		if err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			var raw json.RawMessage = data
			output.PrintJSON(raw)
			return nil
		}

		var resp struct {
			TotalCount int `json:"totalCount"`
			Items      []struct {
				UserName    string  `json:"userName"`
				DisplayName string  `json:"displayName"`
				Gain        float64 `json:"gain"`
				RiskScore   int     `json:"riskScore"`
				Copiers     int     `json:"copiers"`
				AUM         float64 `json:"aum"`
				IsPI        bool    `json:"isPopularInvestor"`
			} `json:"items"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			output.PrintJSON(json.RawMessage(data))
			return nil
		}

		rows := make([]output.UserRow, len(resp.Items))
		for i, u := range resp.Items {
			rows[i] = output.UserRow{
				Username:    u.UserName,
				DisplayName: u.DisplayName,
				Gain:        u.Gain,
				RiskScore:   u.RiskScore,
				Copiers:     u.Copiers,
				AUM:         u.AUM,
				IsPI:        u.IsPI,
			}
		}

		output.PrintUserSearch(rows, output.GetFormat())
		return nil
	},
}

var copyPerformanceCmd = &cobra.Command{
	Use:   "performance <username>",
	Short: "View a trader's historical performance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetUserGain(args[0])
		if err != nil {
			return err
		}

		monthly := make([]output.GainRow, len(resp.Monthly))
		for i, g := range resp.Monthly {
			monthly[i] = output.GainRow{Period: g.Timestamp, Gain: g.Gain}
		}

		yearly := make([]output.GainRow, len(resp.Yearly))
		for i, g := range resp.Yearly {
			yearly[i] = output.GainRow{Period: g.Timestamp, Gain: g.Gain}
		}

		output.PrintUserGain(args[0], monthly, yearly, output.GetFormat())
		return nil
	},
}

var copyCopierCmd = &cobra.Command{
	Use:   "copiers",
	Short: "View who is copying you",
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

func init() {
	copyDiscoverCmd.Flags().String("period", "LastYear", "analysis period (CurrMonth, CurrQuarter, CurrYear, LastYear, LastTwoYears)")
	copyDiscoverCmd.Flags().Int("page", 1, "page number")
	copyDiscoverCmd.Flags().Int("page-size", 20, "results per page")

	copyCmd.AddCommand(copyDiscoverCmd)
	copyCmd.AddCommand(copyPerformanceCmd)
	copyCmd.AddCommand(copyCopierCmd)
	rootCmd.AddCommand(copyCmd)
}
