package output

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type PortfolioSummary struct {
	Credit        float64
	UnrealizedPnL float64
	Equity        float64
	PositionCount int
	OrderCount    int
	CopyCount     int
	IsDemo        bool
}

func PrintPortfolioSummary(s PortfolioSummary, format Format) {
	if format == JSON {
		PrintJSON(s)
		return
	}

	label := "Real"
	if s.IsDemo {
		label = Yellow("Demo")
	}

	t := NewDetailTable()
	DetailRow(t, "Account", label)
	DetailRow(t, "Available Cash", FormatMoney(s.Credit))
	DetailRow(t, "Unrealized P&L", FormatPnL(s.UnrealizedPnL))
	DetailRow(t, "Equity", FormatMoney(s.Equity))
	DetailRow(t, "Open Positions", s.PositionCount)
	DetailRow(t, "Pending Orders", s.OrderCount)
	DetailRow(t, "Active Copies", s.CopyCount)
	RenderTable(t)
}

type PositionRow struct {
	PositionID int
	Symbol     string
	Direction  string
	Amount     float64
	Units      float64
	OpenRate   float64
	Leverage   int
	SL         float64
	TP         float64
	OpenDate   string
}

func PrintPositions(positions []PositionRow, format Format) {
	if format == JSON {
		PrintJSON(positions)
		return
	}

	if len(positions) == 0 {
		Infof("No open positions.")
		return
	}

	t := NewTable("ID", "Symbol", "Side", "Amount", "Units", "Open Rate", "Leverage", "SL", "TP", "Opened")
	for _, p := range positions {
		dir := Green("LONG")
		if p.Direction == "Short" {
			dir = Red("SHORT")
		}
		sl := "-"
		if p.SL > 0 {
			sl = fmt.Sprintf("%.4f", p.SL)
		}
		tp := "-"
		if p.TP > 0 {
			tp = fmt.Sprintf("%.4f", p.TP)
		}
		t.AppendRow(table.Row{
			p.PositionID,
			Cyan(p.Symbol),
			dir,
			FormatMoney(p.Amount),
			fmt.Sprintf("%.4f", p.Units),
			fmt.Sprintf("%.4f", p.OpenRate),
			fmt.Sprintf("%dx", p.Leverage),
			sl,
			tp,
			p.OpenDate,
		})
	}
	RenderTable(t)
}

type OrderRow struct {
	OrderID    int
	Symbol     string
	Direction  string
	Amount     float64
	Units      float64
	Rate       float64
	Leverage   int
	SL         float64
	TP         float64
	OpenDate   string
}

func PrintOrders(orders []OrderRow, format Format) {
	if format == JSON {
		PrintJSON(orders)
		return
	}

	if len(orders) == 0 {
		Infof("No pending orders.")
		return
	}

	t := NewTable("ID", "Symbol", "Side", "Amount", "Rate", "Leverage", "SL", "TP", "Created")
	for _, o := range orders {
		dir := Green("LONG")
		if o.Direction == "Short" {
			dir = Red("SHORT")
		}
		sl := "-"
		if o.SL > 0 {
			sl = fmt.Sprintf("%.4f", o.SL)
		}
		tp := "-"
		if o.TP > 0 {
			tp = fmt.Sprintf("%.4f", o.TP)
		}
		t.AppendRow(table.Row{
			o.OrderID,
			Cyan(o.Symbol),
			dir,
			FormatMoney(o.Amount),
			fmt.Sprintf("%.4f", o.Rate),
			fmt.Sprintf("%dx", o.Leverage),
			sl,
			tp,
			o.OpenDate,
		})
	}
	RenderTable(t)
}

type HistoryRow struct {
	PositionID int64
	Symbol     string
	Direction  string
	Investment float64
	NetProfit  float64
	OpenRate   float64
	CloseRate  float64
	Leverage   int
	Units      float64
	Fees       float64
	Opened     string
	Closed     string
}

func PrintTradeHistory(entries []HistoryRow, format Format) {
	if format == JSON {
		PrintJSON(entries)
		return
	}

	if len(entries) == 0 {
		Infof("No trade history found.")
		return
	}

	t := NewTable("ID", "Symbol", "Side", "Investment", "P&L", "Open", "Close", "Leverage", "Fees", "Closed")
	for _, e := range entries {
		dir := Green("LONG")
		if e.Direction == "Short" {
			dir = Red("SHORT")
		}
		t.AppendRow(table.Row{
			e.PositionID,
			Cyan(e.Symbol),
			dir,
			FormatMoney(e.Investment),
			FormatPnL(e.NetProfit),
			fmt.Sprintf("%.4f", e.OpenRate),
			fmt.Sprintf("%.4f", e.CloseRate),
			fmt.Sprintf("%dx", e.Leverage),
			fmt.Sprintf("$%.2f", e.Fees),
			e.Closed,
		})
	}
	RenderTable(t)
}
