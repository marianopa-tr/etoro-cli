package output

import (
	"encoding/json"
	"testing"
)

func TestPrintWatchlistsTable(t *testing.T) {
	rows := []WatchlistRow{
		{ID: "abc", Name: "Tech Stocks", Type: "Custom", Items: 5, IsDefault: true, Rank: 1},
		{ID: "def", Name: "Crypto", Type: "Custom", Items: 3, Rank: 2},
	}

	out := captureStdout(t, func() {
		PrintWatchlists(rows, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintWatchlistsJSON(t *testing.T) {
	rows := []WatchlistRow{{ID: "abc", Name: "Tech", Items: 5}}

	out := captureStdout(t, func() {
		PrintWatchlists(rows, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintWatchlistsEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintWatchlists(nil, Table)
	})
	_ = out
}

func TestPrintWatchlistItemsTable(t *testing.T) {
	items := []WatchlistItemRow{
		{ItemID: 1001, Symbol: "AAPL", Name: "Apple", Price: 150.25, DailyChg: 1.5},
	}

	out := captureStdout(t, func() {
		PrintWatchlistItems("My List", items, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintWatchlistItemsJSON(t *testing.T) {
	items := []WatchlistItemRow{{ItemID: 1001, Symbol: "AAPL"}}

	out := captureStdout(t, func() {
		PrintWatchlistItems("My List", items, JSON)
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["watchlist"] != "My List" {
		t.Errorf("watchlist = %v", result["watchlist"])
	}
}

func TestPrintCuratedListsTable(t *testing.T) {
	lists := []CuratedListRow{
		{Name: "Top Tech", Description: "Top tech stocks", ItemCount: 10},
	}

	out := captureStdout(t, func() {
		PrintCuratedLists(lists, Table)
	})
	if out == "" {
		t.Error("empty output")
	}
}

func TestPrintCuratedListsJSON(t *testing.T) {
	lists := []CuratedListRow{{Name: "Top Tech", ItemCount: 10}}

	out := captureStdout(t, func() {
		PrintCuratedLists(lists, JSON)
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPrintCuratedListsEmpty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintCuratedLists(nil, Table)
	})
	_ = out
}
