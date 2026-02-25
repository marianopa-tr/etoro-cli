package output

import (
	"encoding/json"
	"testing"
)

func TestPrintPortfolioSummaryTable(t *testing.T) {
	summary := PortfolioSummary{
		Credit: 10000, UnrealizedPnL: 500, Equity: 10500,
		PositionCount: 5, OrderCount: 2, CopyCount: 1,
	}

	out := captureStdout(t, func() {
		PrintPortfolioSummary(summary, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintPortfolioSummaryJSON(t *testing.T) {
	summary := PortfolioSummary{Credit: 10000, Equity: 10500}

	out := captureStdout(t, func() {
		PrintPortfolioSummary(summary, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintPositionsTable(t *testing.T) {
	rows := []PositionRow{
		{PositionID: 1, Symbol: "AAPL", Direction: "Long", Amount: 500, Units: 3.5, OpenRate: 150.0, Leverage: 1},
		{PositionID: 2, Symbol: "TSLA", Direction: "Short", Amount: 300, Units: 1.5, OpenRate: 200.0, Leverage: 2, SL: 210, TP: 180},
	}

	out := captureStdout(t, func() {
		PrintPositions(rows, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintPositionsJSON(t *testing.T) {
	rows := []PositionRow{{PositionID: 1, Symbol: "AAPL", Amount: 500}}

	out := captureStdout(t, func() {
		PrintPositions(rows, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintPositionsEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintPositions(nil, Table)
	})
	_ = out
}

func TestPrintOrdersTable(t *testing.T) {
	rows := []OrderRow{
		{OrderID: 10, Symbol: "AAPL", Direction: "Long", Amount: 500, Rate: 150, Leverage: 1},
	}

	out := captureStdout(t, func() {
		PrintOrders(rows, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintOrdersJSON(t *testing.T) {
	rows := []OrderRow{{OrderID: 10, Symbol: "AAPL", Amount: 500}}

	out := captureStdout(t, func() {
		PrintOrders(rows, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintOrdersEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintOrders(nil, Table)
	})
	_ = out
}

func TestPrintTradeHistoryTable(t *testing.T) {
	rows := []HistoryRow{
		{PositionID: 100, Symbol: "AAPL", Direction: "Long", Investment: 500, NetProfit: 50, OpenRate: 150, CloseRate: 160, Leverage: 1},
	}

	out := captureStdout(t, func() {
		PrintTradeHistory(rows, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintTradeHistoryJSON(t *testing.T) {
	rows := []HistoryRow{{PositionID: 100, Symbol: "AAPL", NetProfit: 50}}

	out := captureStdout(t, func() {
		PrintTradeHistory(rows, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintTradeHistoryEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintTradeHistory(nil, Table)
	})
	_ = out
}
