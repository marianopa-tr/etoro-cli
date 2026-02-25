package output

import (
	"encoding/json"
	"testing"
)

func TestPrintUserSearchTable(t *testing.T) {
	users := []UserRow{
		{Username: "trader1", DisplayName: "Top Trader", Gain: 25.5, RiskScore: 3, Copiers: 100, AUM: 500000, IsPI: true},
		{Username: "trader2", DisplayName: "Risky Guy", Gain: -5.0, RiskScore: 8, Copiers: 10, AUM: 50000},
	}

	out := captureStdout(t, func() {
		PrintUserSearch(users, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintUserSearchJSON(t *testing.T) {
	users := []UserRow{{Username: "trader1", Gain: 25.5}}

	out := captureStdout(t, func() {
		PrintUserSearch(users, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintUserSearchEmpty(t *testing.T) {
	captureStdout(t, func() {
		PrintUserSearch(nil, Table)
	})
}

func TestRiskLabel(t *testing.T) {
	tests := []struct {
		score int
	}{
		{1}, {3}, {5}, {6}, {7}, {10},
	}
	for _, tt := range tests {
		got := riskLabel(tt.score)
		if got == "" {
			t.Errorf("riskLabel(%d) returned empty", tt.score)
		}
	}
}

func TestPrintUserGainTable(t *testing.T) {
	monthly := []GainRow{{Period: "2025-01", Gain: 5.5}, {Period: "2025-02", Gain: -2.0}}
	yearly := []GainRow{{Period: "2024", Gain: 25.0}}

	out := captureStdout(t, func() {
		PrintUserGain("trader1", monthly, yearly, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintUserGainJSON(t *testing.T) {
	monthly := []GainRow{{Period: "2025-01", Gain: 5.5}}
	yearly := []GainRow{{Period: "2024", Gain: 25.0}}

	out := captureStdout(t, func() {
		PrintUserGain("trader1", monthly, yearly, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["username"] != "trader1" {
		t.Errorf("username = %v", result["username"])
	}
}

func TestPrintCopiersTable(t *testing.T) {
	copiers := []CopierRow{
		{Country: "US", Club: "Gold", Duration: "1-3 months", Amount: "$1K-$5K"},
	}

	out := captureStdout(t, func() {
		PrintCopiers(copiers, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintCopiersJSON(t *testing.T) {
	copiers := []CopierRow{{Country: "US"}}

	out := captureStdout(t, func() {
		PrintCopiers(copiers, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintCopiersEmpty(t *testing.T) {
	captureStdout(t, func() {
		PrintCopiers(nil, Table)
	})
}
