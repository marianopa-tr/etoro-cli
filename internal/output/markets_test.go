package output

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestPrintInstrumentSearchTable(t *testing.T) {
	SetFormat(Table)
	rows := []InstrumentRow{
		{ID: 1001, Symbol: "AAPL", Name: "Apple Inc", Type: "Stocks", Price: 150.25, DailyChg: 1.5, Open: true},
		{ID: 1002, Symbol: "TSLA", Name: "Tesla Inc", Type: "Stocks", Price: 200.50, DailyChg: -2.3, Open: false},
	}

	out := captureStdout(t, func() {
		PrintInstrumentSearch(rows, 100, Table)
	})

	if out == "" {
		t.Error("output is empty")
	}
	if len(out) < 20 {
		t.Errorf("output seems too short: %q", out)
	}
}

func TestPrintInstrumentSearchJSON(t *testing.T) {
	rows := []InstrumentRow{
		{ID: 1001, Symbol: "AAPL", Name: "Apple Inc", Price: 150.25},
	}

	out := captureStdout(t, func() {
		PrintInstrumentSearch(rows, 1, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if result["totalItems"] != float64(1) {
		t.Errorf("totalItems = %v", result["totalItems"])
	}
}

func TestPrintInstrumentSearchEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintInstrumentSearch(nil, 0, Table)
	})
	_ = out
}

func TestPrintInstrumentDetailTable(t *testing.T) {
	detail := InstrumentDetail{
		ID: 1001, Symbol: "AAPL", Name: "Apple Inc", Type: "Stocks",
		Price: 150.25, DailyChg: 1.5, WeeklyChg: 3.2,
		Open: true, Tradable: true, BuyEnabled: true,
	}

	out := captureStdout(t, func() {
		PrintInstrumentDetail(detail, Table)
	})

	if out == "" {
		t.Error("output is empty")
	}
}

func TestPrintInstrumentDetailJSON(t *testing.T) {
	detail := InstrumentDetail{
		ID: 1001, Symbol: "AAPL", Name: "Apple Inc",
	}

	out := captureStdout(t, func() {
		PrintInstrumentDetail(detail, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintQuotesTable(t *testing.T) {
	quotes := []QuoteRow{
		{Symbol: "AAPL", Name: "Apple", Bid: 150.20, Ask: 150.25, Last: 150.22},
	}

	out := captureStdout(t, func() {
		PrintQuotes(quotes, Table)
	})

	if out == "" {
		t.Error("output is empty")
	}
}

func TestPrintQuotesJSON(t *testing.T) {
	quotes := []QuoteRow{
		{Symbol: "AAPL", Bid: 150.20, Ask: 150.25},
	}

	out := captureStdout(t, func() {
		PrintQuotes(quotes, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("count = %d", len(result))
	}
}

func TestPrintQuotesEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintQuotes(nil, Table)
	})
	_ = out
}
