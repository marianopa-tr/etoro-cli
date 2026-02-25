package output

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type UserRow struct {
	Username    string
	DisplayName string
	Gain        float64
	RiskScore   int
	Copiers     int
	AUM         float64
	IsPI        bool
}

func PrintUserSearch(users []UserRow, format Format) {
	if format == JSON {
		PrintJSON(users)
		return
	}

	if len(users) == 0 {
		Infof("No users found.")
		return
	}

	t := NewTable("Username", "Name", "Gain", "Risk", "Copiers", "AUM", "PI")
	for _, u := range users {
		pi := ""
		if u.IsPI {
			pi = Green("★")
		}
		t.AppendRow(table.Row{
			Cyan(u.Username),
			u.DisplayName,
			FormatPercent(u.Gain),
			riskLabel(u.RiskScore),
			u.Copiers,
			FormatMoney(u.AUM),
			pi,
		})
	}
	RenderTable(t)
}

func riskLabel(score int) string {
	switch {
	case score <= 3:
		return Green(fmt.Sprintf("%d", score))
	case score <= 6:
		return Yellow(fmt.Sprintf("%d", score))
	default:
		return Red(fmt.Sprintf("%d", score))
	}
}

type TraderProfileRow struct {
	Username           string
	FullName           string
	IsPopularInvestor  bool
	Gain               float64
	DailyGain          float64
	WeeklyGain         float64
	RiskScore          int
	MaxDailyRisk       int
	MaxMonthlyRisk     int
	Copiers            int
	Trades             int
	WinRatio           float64
	DailyDD            float64
	WeeklyDD           float64
	PeakToValley       float64
	ProfitWeeksPct     float64
	ProfitMonthsPct    float64
	AvgPosSize         float64
	AUMTierDesc        string
	WeeksRegistered    int
}

func PrintTraderProfile(p TraderProfileRow, format Format) {
	if format == JSON {
		PrintJSON(p)
		return
	}

	pi := ""
	if p.IsPopularInvestor {
		pi = " " + Green("★ Popular Investor")
	}

	t := NewDetailTable()
	DetailRow(t, "Username", Cyan(p.Username)+pi)
	DetailRow(t, "Full Name", p.FullName)
	DetailRow(t, "AUM", p.AUMTierDesc)
	DetailRow(t, "Weeks Active", p.WeeksRegistered)
	DetailRow(t, "Copiers", p.Copiers)
	DetailRow(t, "Gain (period)", FormatPercent(p.Gain))
	DetailRow(t, "Daily Gain", FormatPercent(p.DailyGain))
	DetailRow(t, "Weekly Gain", FormatPercent(p.WeeklyGain))
	DetailRow(t, "Risk Score", riskLabel(p.RiskScore))
	DetailRow(t, "Max Daily Risk", riskLabel(p.MaxDailyRisk))
	DetailRow(t, "Max Monthly Risk", riskLabel(p.MaxMonthlyRisk))
	DetailRow(t, "Trades", p.Trades)
	DetailRow(t, "Win Ratio", FormatPercent(p.WinRatio))
	DetailRow(t, "Daily Drawdown", FormatPercent(p.DailyDD))
	DetailRow(t, "Weekly Drawdown", FormatPercent(p.WeeklyDD))
	DetailRow(t, "Peak to Valley", FormatPercent(p.PeakToValley))
	DetailRow(t, "Profitable Weeks", FormatPercent(p.ProfitWeeksPct))
	DetailRow(t, "Profitable Months", FormatPercent(p.ProfitMonthsPct))
	DetailRow(t, "Avg Position Size", FormatPercent(p.AvgPosSize))
	RenderTable(t)
}

type GainRow struct {
	Period string
	Gain   float64
}

func PrintUserGain(username string, monthly []GainRow, yearly []GainRow, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"username": username,
			"monthly":  monthly,
			"yearly":   yearly,
		})
		return
	}

	fmt.Printf("\n  Performance: %s\n\n", Bold(username))

	if len(yearly) > 0 {
		fmt.Println("  " + Bold("Yearly"))
		t := NewTable("Period", "Gain")
		for _, g := range yearly {
			t.AppendRow(table.Row{g.Period, FormatPercent(g.Gain)})
		}
		RenderTable(t)
		fmt.Println()
	}

	if len(monthly) > 0 {
		fmt.Println("  " + Bold("Monthly"))
		t := NewTable("Period", "Gain")
		for _, g := range monthly {
			t.AppendRow(table.Row{g.Period, FormatPercent(g.Gain)})
		}
		RenderTable(t)
	}
}

type CopierRow struct {
	Country     string
	Club        string
	Duration    string
	Amount      string
}

func PrintCopiers(copiers []CopierRow, format Format) {
	if format == JSON {
		PrintJSON(copiers)
		return
	}

	if len(copiers) == 0 {
		Infof("No copiers found.")
		return
	}

	t := NewTable("Country", "Club", "Duration", "Amount")
	for _, c := range copiers {
		t.AppendRow(table.Row{c.Country, c.Club, c.Duration, c.Amount})
	}
	RenderTable(t)
}
