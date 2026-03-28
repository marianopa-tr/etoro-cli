package cmd

import (
	"fmt"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/marianopa-tr/etoro-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "View and create social feed posts",
	Long: `Interact with eToro's social feed: view discussions about instruments
or users, and create new posts.

Examples:
  etoro feed list --instrument AAPL
  etoro feed list --username BigTech
  etoro feed post "Bullish on AAPL this quarter!"`,
}

var feedListCmd = &cobra.Command{
	Use:   "list",
	Short: "View feed posts",
	Long: `View feed posts for an instrument or user.

Examples:
  etoro feed list --instrument AAPL
  etoro feed list --instrument BTC --limit 20
  etoro feed list --username BigTech`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		instrument, _ := cmd.Flags().GetString("instrument")
		username, _ := cmd.Flags().GetString("username")
		limit, _ := cmd.Flags().GetInt("limit")

		var resp *api.FeedResponse
		var err error

		if instrument != "" {
			res := resolver.New(client)
			instrumentID, _, resolveErr := res.Resolve(instrument)
			if resolveErr != nil {
				return resolveErr
			}
			resp, err = client.GetInstrumentFeed(instrumentID, 0, limit)
		} else if username != "" {
			profile, profileErr := client.GetUserProfile(username)
			if profileErr != nil {
				return profileErr
			}
			resp, err = client.GetUserFeed(fmt.Sprintf("%d", profile.GCID), 0, limit)
		} else {
			return errorf("specify --instrument <symbol> or --username <name>")
		}

		if err != nil {
			return err
		}

		rows := make([]output.FeedPostRow, 0, len(resp.Discussions))
		for _, d := range resp.Discussions {
			rows = append(rows, output.FeedPostRow{
				ID:      d.ID,
				Author:  d.Post.Owner.Username,
				Message: d.Post.Message.Text,
				Created: d.Post.Created,
			})
		}

		output.PrintFeedPosts(rows, output.GetFormat())
		return nil
	},
}

var feedPostCmd = &cobra.Command{
	Use:   "post <message>",
	Short: "Create a new discussion post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		instrument, _ := cmd.Flags().GetString("instrument")

		req := &api.CreatePostRequest{
			Message: args[0],
		}

		if instrument != "" {
			res := resolver.New(client)
			instrumentID, _, err := res.Resolve(instrument)
			if err != nil {
				return err
			}
			req.Tags = &api.PostTagsRequest{
				Instruments: []api.PostTagInstrument{{ID: instrumentID}},
			}
		}

		_, err := client.CreatePost(req)
		if err != nil {
			return err
		}

		output.PrintPostCreated(output.GetFormat())
		return nil
	},
}

func init() {
	feedListCmd.Flags().String("instrument", "", "filter by instrument symbol")
	feedListCmd.Flags().String("username", "", "filter by username")
	feedListCmd.Flags().Int("limit", 10, "number of posts to show")

	feedPostCmd.Flags().String("instrument", "", "tag an instrument in the post")

	feedCmd.AddCommand(feedListCmd)
	feedCmd.AddCommand(feedPostCmd)
	rootCmd.AddCommand(feedCmd)
}
