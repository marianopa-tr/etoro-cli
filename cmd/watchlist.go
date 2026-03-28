package cmd

import (
	"fmt"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/marianopa-tr/etoro-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var watchlistCmd = &cobra.Command{
	Use:   "watchlist",
	Short: "Manage your watchlists",
	Long: `Create, view, and manage watchlists.

Examples:
  etoro watchlist list
  etoro watchlist create "My Tech Stocks"
  etoro watchlist add AAPL
  etoro watchlist add TSLA --to <watchlistId>
  etoro watchlist remove AAPL`,
}

var watchlistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your watchlists",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		curated, _ := cmd.Flags().GetBool("curated")
		if curated {
			resp, err := client.GetCuratedLists()
			if err != nil {
				return err
			}
			rows := make([]output.CuratedListRow, len(resp.CuratedLists))
			for i, l := range resp.CuratedLists {
				rows[i] = output.CuratedListRow{
					Name:        l.Name,
					Description: l.Description,
					ItemCount:   len(l.Items),
				}
			}
			output.PrintCuratedLists(rows, output.GetFormat())
			return nil
		}

		resp, err := client.GetWatchlists()
		if err != nil {
			return err
		}

		rows := make([]output.WatchlistRow, len(resp.Watchlists))
		for i, w := range resp.Watchlists {
			rows[i] = output.WatchlistRow{
				ID:        w.WatchlistID,
				Name:      w.Name,
				Type:      w.WatchlistType,
				Items:     w.TotalItems,
				IsDefault: w.IsDefault || w.IsUserSelectedDefault,
				Rank:      w.WatchlistRank,
			}
		}

		output.PrintWatchlists(rows, output.GetFormat())
		return nil
	},
}

var watchlistCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		_, err := client.CreateWatchlist(&api.CreateWatchlistRequest{
			Name: args[0],
		})
		if err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "created", "name": args[0]})
		} else {
			output.Successf("Watchlist %q created.", args[0])
		}
		return nil
	},
}

var watchlistAddCmd = &cobra.Command{
	Use:   "add <symbol>",
	Short: "Add an instrument to a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		instrumentID, symbol, err := res.Resolve(args[0])
		if err != nil {
			return err
		}

		watchlistID, _ := cmd.Flags().GetString("to")
		if watchlistID == "" {
			wl, dlErr := defaultWatchlist(client)
			if dlErr != nil {
				return dlErr
			}
			watchlistID = wl
		}

		_, err = client.AddWatchlistItems(watchlistID, []api.WatchlistItem{
			{ItemID: instrumentID, ItemType: "Instrument"},
		})
		if err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "added", "symbol": symbol, "watchlistId": watchlistID})
		} else {
			output.Successf("Added %s to watchlist %s.", output.Cyan(symbol), watchlistID)
		}
		return nil
	},
}

var watchlistRemoveCmd = &cobra.Command{
	Use:   "remove <symbol>",
	Short: "Remove an instrument from a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		instrumentID, symbol, err := res.Resolve(args[0])
		if err != nil {
			return err
		}

		watchlistID, _ := cmd.Flags().GetString("from")
		if watchlistID == "" {
			wl, dlErr := defaultWatchlist(client)
			if dlErr != nil {
				return dlErr
			}
			watchlistID = wl
		}

		err = client.RemoveWatchlistItems(watchlistID, []api.WatchlistItem{
			{ItemID: instrumentID, ItemType: "Instrument"},
		})
		if err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "removed", "symbol": symbol, "watchlistId": watchlistID})
		} else {
			output.Successf("Removed %s from watchlist %s.", output.Cyan(symbol), watchlistID)
		}
		return nil
	},
}

func defaultWatchlist(client *api.Client) (string, error) {
	resp, err := client.GetWatchlists()
	if err != nil {
		return "", fmt.Errorf("fetching watchlists to find default: %w", err)
	}
	for _, w := range resp.Watchlists {
		if w.IsDefault || w.IsUserSelectedDefault {
			return w.WatchlistID, nil
		}
	}
	if len(resp.Watchlists) > 0 {
		return resp.Watchlists[0].WatchlistID, nil
	}
	return "", fmt.Errorf("no watchlists found; create one with: etoro watchlist create <name>")
}

func init() {
	watchlistListCmd.Flags().Bool("curated", false, "show curated/recommended lists instead")

	watchlistAddCmd.Flags().String("to", "", "target watchlist ID")
	watchlistRemoveCmd.Flags().String("from", "", "source watchlist ID")

	watchlistCmd.AddCommand(watchlistListCmd)
	watchlistCmd.AddCommand(watchlistCreateCmd)
	watchlistCmd.AddCommand(watchlistAddCmd)
	watchlistCmd.AddCommand(watchlistRemoveCmd)
	rootCmd.AddCommand(watchlistCmd)
}
