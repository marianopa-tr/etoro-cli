package output

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type InstrumentRow struct {
	ID         int
	Symbol     string
	Name       string
	Type       string
	Exchange   string
	Price      float64
	DailyChg   float64
	WeeklyChg  float64
	Open       bool
	Tradable   bool
	Popularity int
}

func PrintInstrumentSearch(instruments []InstrumentRow, total int, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"totalItems":  total,
			"instruments": instruments,
		})
		return
	}

	if len(instruments) == 0 {
		Infof("No instruments found.")
		return
	}

	t := NewTable("ID", "Symbol", "Name", "Type", "Price", "Daily", "Weekly", "Status")
	for _, inst := range instruments {
		status := FormatBool(inst.Open, "Open", "Closed")
		t.AppendRow(table.Row{
			inst.ID,
			Cyan(inst.Symbol),
			inst.Name,
			inst.Type,
			fmt.Sprintf("%.2f", inst.Price),
			FormatPercent(inst.DailyChg),
			FormatPercent(inst.WeeklyChg),
			status,
		})
	}
	t.AppendFooter(table.Row{"", "", "", "", "", "", "", fmt.Sprintf("%d results", total)})
	RenderTable(t)
}

type InstrumentDetail struct {
	ID            int
	Symbol        string
	Name          string
	Type          string
	Exchange      string
	AssetClass    string
	Price         float64
	DailyChg      float64
	WeeklyChg     float64
	MonthlyChg    float64
	ThreeMonthChg float64
	OneYearChg    float64
	Open          bool
	Tradable      bool
	BuyEnabled    bool
	Popularity    int
}

func PrintInstrumentDetail(inst InstrumentDetail, format Format) {
	if format == JSON {
		PrintJSON(inst)
		return
	}

	t := NewDetailTable()
	DetailRow(t, "Instrument ID", inst.ID)
	DetailRow(t, "Symbol", Cyan(inst.Symbol))
	DetailRow(t, "Name", inst.Name)
	DetailRow(t, "Type", inst.Type)
	DetailRow(t, "Exchange", inst.Exchange)
	DetailRow(t, "Asset Class", inst.AssetClass)
	DetailRow(t, "Current Price", fmt.Sprintf("%.4f", inst.Price))
	DetailRow(t, "Daily Change", FormatPercent(inst.DailyChg))
	DetailRow(t, "Weekly Change", FormatPercent(inst.WeeklyChg))
	DetailRow(t, "Monthly Change", FormatPercent(inst.MonthlyChg))
	DetailRow(t, "3-Month Change", FormatPercent(inst.ThreeMonthChg))
	DetailRow(t, "1-Year Change", FormatPercent(inst.OneYearChg))
	DetailRow(t, "Market Status", FormatBool(inst.Open, "Open", "Closed"))
	DetailRow(t, "Buy Enabled", FormatBool(inst.BuyEnabled, "Yes", "No"))
	DetailRow(t, "Popularity (7d)", inst.Popularity)
	RenderTable(t)
}

type QuoteRow struct {
	Symbol string
	Name   string
	Bid    float64
	Ask    float64
	Last   float64
	Spread float64
}

func PrintQuotes(quotes []QuoteRow, format Format) {
	if format == JSON {
		PrintJSON(quotes)
		return
	}

	if len(quotes) == 0 {
		Infof("No quotes available.")
		return
	}

	t := NewTable("Symbol", "Name", "Bid", "Ask", "Spread", "Last")
	for _, q := range quotes {
		spread := q.Ask - q.Bid
		t.AppendRow(table.Row{
			Cyan(q.Symbol),
			q.Name,
			fmt.Sprintf("%.4f", q.Bid),
			fmt.Sprintf("%.4f", q.Ask),
			fmt.Sprintf("%.4f", spread),
			fmt.Sprintf("%.4f", q.Last),
		})
	}
	RenderTable(t)
}
