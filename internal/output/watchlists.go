package output

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type WatchlistRow struct {
	ID         string
	Name       string
	Type       string
	Items      int
	IsDefault  bool
	Rank       int
}

func PrintWatchlists(watchlists []WatchlistRow, format Format) {
	if format == JSON {
		PrintJSON(watchlists)
		return
	}

	if len(watchlists) == 0 {
		Infof("No watchlists found.")
		return
	}

	t := NewTable("ID", "Name", "Type", "Items", "Default", "Rank")
	for _, w := range watchlists {
		def := ""
		if w.IsDefault {
			def = Green("★")
		}
		t.AppendRow(table.Row{
			w.ID,
			w.Name,
			w.Type,
			w.Items,
			def,
			w.Rank,
		})
	}
	RenderTable(t)
}

type WatchlistItemRow struct {
	ItemID   int
	Symbol   string
	Name     string
	Type     string
	Price    float64
	DailyChg float64
}

func PrintWatchlistItems(name string, items []WatchlistItemRow, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"watchlist": name,
			"items":     items,
		})
		return
	}

	if len(items) == 0 {
		Infof("Watchlist %q is empty.", name)
		return
	}

	fmt.Printf("\n  %s\n\n", Bold(name))
	t := NewTable("ID", "Symbol", "Name", "Type", "Price", "Daily")
	for _, item := range items {
		t.AppendRow(table.Row{
			item.ItemID,
			Cyan(item.Symbol),
			item.Name,
			item.Type,
			fmt.Sprintf("%.2f", item.Price),
			FormatPercent(item.DailyChg),
		})
	}
	RenderTable(t)
}

type CuratedListRow struct {
	Name        string
	Description string
	ItemCount   int
}

func PrintCuratedLists(lists []CuratedListRow, format Format) {
	if format == JSON {
		PrintJSON(lists)
		return
	}

	if len(lists) == 0 {
		Infof("No curated lists found.")
		return
	}

	t := NewTable("Name", "Description", "Items")
	for _, l := range lists {
		t.AppendRow(table.Row{
			Bold(l.Name),
			l.Description,
			l.ItemCount,
		})
	}
	RenderTable(t)
}
