package cmd

import (
	"fmt"
	"strings"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var agentPortfolioCmd = &cobra.Command{
	Use:     "agent-portfolio",
	Aliases: []string{"ap"},
	Short:   "Manage agent portfolios and their user tokens",
	Long: `Create, list, inspect, and delete agent portfolios.
Agent portfolios are dedicated sub-accounts with their own virtual balance
that mirror trades proportionally into your real account.

This feature requires real account keys (not demo).

Examples:
  etoro agent-portfolio list
  etoro agent-portfolio create --name MyBot01 --investment 2000 --token-name bot-token
  etoro agent-portfolio get <portfolioId>
  etoro agent-portfolio delete <portfolioId>
  etoro agent-portfolio token create <portfolioId> --name read-token --scopes 200
  etoro agent-portfolio token delete <portfolioId> <tokenId>
  etoro agent-portfolio token update <portfolioId> <tokenId> --scopes 200,202`,
}

var apListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agent portfolios",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, false)

		resp, err := client.GetAgentPortfolios()
		if err != nil {
			return err
		}

		rows := make([]output.AgentPortfolioRow, len(resp.AgentPortfolios))
		for i, p := range resp.AgentPortfolios {
			rows[i] = output.AgentPortfolioRow{
				ID:             p.AgentPortfolioID,
				Name:           p.AgentPortfolioName,
				GCID:           p.AgentPortfolioGCID,
				VirtualBalance: p.AgentPortfolioVirtualBalance,
				MirrorID:       p.MirrorID,
				CreatedAt:      p.CreatedAt,
				TokenCount:     len(p.UserTokens),
			}
		}
		output.PrintAgentPortfolios(rows, output.GetFormat())
		return nil
	},
}

var apGetCmd = &cobra.Command{
	Use:   "get <portfolioId>",
	Short: "Show details of an agent portfolio",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, false)

		resp, err := client.GetAgentPortfolios()
		if err != nil {
			return err
		}

		var found *api.AgentPortfolioItem
		for i, p := range resp.AgentPortfolios {
			if p.AgentPortfolioID == args[0] || p.AgentPortfolioName == args[0] {
				found = &resp.AgentPortfolios[i]
				break
			}
		}
		if found == nil {
			return errorf("agent portfolio %q not found", args[0])
		}

		tokens := make([]output.AgentPortfolioTokenRow, len(found.UserTokens))
		for i, tok := range found.UserTokens {
			tokens[i] = output.AgentPortfolioTokenRow{
				TokenID:   tok.UserTokenID,
				Name:      tok.UserTokenName,
				ClientID:  tok.ClientID,
				AppName:   tok.ExternalApplicationName,
				Scopes:    output.FormatScopeIDs(tok.ScopeIDs),
				IPs:       strings.Join(tok.IPsWhitelist, ", "),
				ExpiresAt: tok.ExpiresAt,
				CreatedAt: tok.CreatedAt,
			}
		}

		output.PrintAgentPortfolioDetail(output.AgentPortfolioDetail{
			ID:             found.AgentPortfolioID,
			Name:           found.AgentPortfolioName,
			GCID:           found.AgentPortfolioGCID,
			VirtualBalance: found.AgentPortfolioVirtualBalance,
			MirrorID:       found.MirrorID,
			CreatedAt:      found.CreatedAt,
			Tokens:         tokens,
		}, output.GetFormat())
		return nil
	},
}

var apCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new agent portfolio",
	Long: `Create a new agent portfolio with a dedicated virtual balance.
The --investment amount is deducted from YOUR account to copy-trade
this portfolio proportionally.

Scope IDs: 200=real:read, 201=demo:read, 202=real:write, 203=demo:write

Examples:
  etoro agent-portfolio create --name MyBot01 --investment 2000 --token-name bot-key
  etoro agent-portfolio create --name MyBot01 --investment 500 --token-name tk1 --scopes 200,202
  etoro agent-portfolio create --name MyBot01 --investment 1000 --token-name tk1 --description "Momentum strategy"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		investment, _ := cmd.Flags().GetFloat64("investment")
		tokenName, _ := cmd.Flags().GetString("token-name")
		description, _ := cmd.Flags().GetString("description")
		scopesStr, _ := cmd.Flags().GetString("scopes")
		ipsStr, _ := cmd.Flags().GetString("ips")
		expiresAt, _ := cmd.Flags().GetString("expires")

		if name == "" {
			return errorf("--name is required (6-10 characters)")
		}
		if investment <= 0 {
			return errorf("--investment is required (amount in USD)")
		}
		if tokenName == "" {
			return errorf("--token-name is required")
		}

		scopeIDs := parseScopeIDs(scopesStr)
		if len(scopeIDs) == 0 {
			scopeIDs = []int{200, 202}
		}

		var ips []string
		if ipsStr != "" {
			ips = strings.Split(ipsStr, ",")
		}

		if !flagYes {
			msg := fmt.Sprintf("Create agent portfolio %q with $%.2f investment from your real account?", name, investment)
			if !output.Confirm(msg) {
				output.Infof("Cancelled.")
				return nil
			}
		}

		client := api.NewClient(cfg, false)
		req := &api.CreateAgentPortfolioRequest{
			InvestmentAmountInUsd:     investment,
			AgentPortfolioName:        name,
			AgentPortfolioDescription: description,
			UserTokenName:             tokenName,
			ScopeIDs:                  scopeIDs,
			IPsWhitelist:              ips,
			ExpiresAt:                 expiresAt,
		}

		resp, err := client.CreateAgentPortfolio(req)
		if err != nil {
			return err
		}

		tokens := make([]output.CreatedTokenDisplay, len(resp.UserTokens))
		for i, tok := range resp.UserTokens {
			tokens[i] = output.CreatedTokenDisplay{
				TokenID:     tok.UserTokenID,
				TokenName:   tok.UserTokenName,
				TokenSecret: tok.UserToken,
				Scopes:      output.FormatScopeIDs(tok.ScopeIDs),
			}
		}

		output.PrintCreatedAgentPortfolio(
			resp.AgentPortfolioID,
			resp.AgentPortfolioName,
			resp.AgentPortfolioGCID,
			resp.AgentPortfolioVirtualBalance,
			resp.MirrorID,
			tokens,
			output.GetFormat(),
		)
		return nil
	},
}

var apDeleteCmd = &cobra.Command{
	Use:   "delete <portfolioId>",
	Short: "Delete an agent portfolio",
	Long: `Permanently delete an agent portfolio. This revokes all user tokens,
stops the copy mirror, and removes the portfolio from storage.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		if !flagYes {
			if !output.ConfirmDanger(fmt.Sprintf("This will permanently delete agent portfolio %s and revoke all its tokens.", args[0])) {
				output.Infof("Cancelled.")
				return nil
			}
		}

		client := api.NewClient(cfg, false)
		if err := client.DeleteAgentPortfolio(args[0]); err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "deleted", "portfolioId": args[0]})
		} else {
			output.Successf("Agent portfolio %s deleted.", args[0])
		}
		return nil
	},
}

// --- Token subcommands ---

var apTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage user tokens for an agent portfolio",
}

var apTokenCreateCmd = &cobra.Command{
	Use:   "create <portfolioId>",
	Short: "Create a new user token for an agent portfolio",
	Long: `Create a new user token with specific scopes for an agent portfolio.

Scope IDs: 200=real:read, 201=demo:read, 202=real:write, 203=demo:write

Examples:
  etoro agent-portfolio token create <portfolioId> --name read-token --scopes 200
  etoro agent-portfolio token create <portfolioId> --name full-token --scopes 200,202 --ips 1.2.3.4`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		tokenName, _ := cmd.Flags().GetString("name")
		scopesStr, _ := cmd.Flags().GetString("scopes")
		ipsStr, _ := cmd.Flags().GetString("ips")
		expiresAt, _ := cmd.Flags().GetString("expires")

		if tokenName == "" {
			return errorf("--name is required for the token")
		}

		scopeIDs := parseScopeIDs(scopesStr)
		if len(scopeIDs) == 0 {
			scopeIDs = []int{200, 202}
		}

		var ips []string
		if ipsStr != "" {
			ips = strings.Split(ipsStr, ",")
		}

		client := api.NewClient(cfg, false)
		req := &api.CreateUserTokenRequest{
			UserTokenName: tokenName,
			ScopeIDs:      scopeIDs,
			IPsWhitelist:  ips,
			ExpiresAt:     expiresAt,
		}

		resp, err := client.CreateUserToken(args[0], req)
		if err != nil {
			return err
		}

		output.PrintCreatedUserToken(resp.UserTokenID, resp.UserToken, output.GetFormat())
		return nil
	},
}

var apTokenUpdateCmd = &cobra.Command{
	Use:   "update <portfolioId> <tokenId>",
	Short: "Update a user token's settings",
	Long: `Update the scopes, IP whitelist, or expiration of a user token.

Examples:
  etoro agent-portfolio token update <portfolioId> <tokenId> --scopes 200,201,202
  etoro agent-portfolio token update <portfolioId> <tokenId> --ips 10.0.0.1,10.0.0.2
  etoro agent-portfolio token update <portfolioId> <tokenId> --expires 2027-12-31T23:59:59Z`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		scopesStr, _ := cmd.Flags().GetString("scopes")
		ipsStr, _ := cmd.Flags().GetString("ips")
		expiresAt, _ := cmd.Flags().GetString("expires")

		req := &api.UpdateUserTokenRequest{}
		hasChange := false

		if scopesStr != "" {
			req.ScopeIDs = parseScopeIDs(scopesStr)
			hasChange = true
		}
		if ipsStr != "" {
			req.IPsWhitelist = strings.Split(ipsStr, ",")
			hasChange = true
		}
		if expiresAt != "" {
			req.ExpiresAt = expiresAt
			hasChange = true
		}

		if !hasChange {
			return errorf("at least one of --scopes, --ips, or --expires must be provided")
		}

		client := api.NewClient(cfg, false)
		if err := client.UpdateUserToken(args[0], args[1], req); err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "updated", "portfolioId": args[0], "tokenId": args[1]})
		} else {
			output.Successf("User token %s updated.", args[1])
		}
		return nil
	},
}

var apTokenDeleteCmd = &cobra.Command{
	Use:   "delete <portfolioId> <tokenId>",
	Short: "Revoke and delete a user token",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		if !flagYes {
			if !output.ConfirmDanger(fmt.Sprintf("Permanently revoke user token %s?", args[1])) {
				output.Infof("Cancelled.")
				return nil
			}
		}

		client := api.NewClient(cfg, false)
		if err := client.DeleteUserToken(args[0], args[1]); err != nil {
			return err
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{"status": "deleted", "portfolioId": args[0], "tokenId": args[1]})
		} else {
			output.Successf("User token %s revoked.", args[1])
		}
		return nil
	},
}

func parseScopeIDs(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := parseInt(p)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func init() {
	apCreateCmd.Flags().String("name", "", "portfolio name (6-10 characters)")
	apCreateCmd.Flags().Float64("investment", 0, "investment amount in USD (deducted from your account)")
	apCreateCmd.Flags().String("token-name", "", "name for the initial user token")
	apCreateCmd.Flags().String("description", "", "portfolio description/strategy")
	apCreateCmd.Flags().String("scopes", "200,202", "comma-separated scope IDs (200=real:read, 201=demo:read, 202=real:write, 203=demo:write)")
	apCreateCmd.Flags().String("ips", "", "comma-separated IP whitelist")
	apCreateCmd.Flags().String("expires", "", "token expiration (e.g. 2027-12-31T23:59:59Z)")

	apTokenCreateCmd.Flags().String("name", "", "token name")
	apTokenCreateCmd.Flags().String("scopes", "200,202", "comma-separated scope IDs")
	apTokenCreateCmd.Flags().String("ips", "", "comma-separated IP whitelist")
	apTokenCreateCmd.Flags().String("expires", "", "token expiration (e.g. 2027-12-31T23:59:59Z)")

	apTokenUpdateCmd.Flags().String("scopes", "", "updated scope IDs")
	apTokenUpdateCmd.Flags().String("ips", "", "updated IP whitelist")
	apTokenUpdateCmd.Flags().String("expires", "", "updated expiration")

	apTokenCmd.AddCommand(apTokenCreateCmd)
	apTokenCmd.AddCommand(apTokenUpdateCmd)
	apTokenCmd.AddCommand(apTokenDeleteCmd)

	agentPortfolioCmd.AddCommand(apListCmd)
	agentPortfolioCmd.AddCommand(apGetCmd)
	agentPortfolioCmd.AddCommand(apCreateCmd)
	agentPortfolioCmd.AddCommand(apDeleteCmd)
	agentPortfolioCmd.AddCommand(apTokenCmd)

	rootCmd.AddCommand(agentPortfolioCmd)
}
