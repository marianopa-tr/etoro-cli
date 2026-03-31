package output

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type AgentPortfolioRow struct {
	ID             string
	Name           string
	GCID           int
	VirtualBalance float64
	MirrorID       int
	CreatedAt      string
	TokenCount     int
}

func PrintAgentPortfolios(portfolios []AgentPortfolioRow, format Format) {
	if format == JSON {
		PrintJSON(portfolios)
		return
	}

	if len(portfolios) == 0 {
		Infof("No agent portfolios found.")
		return
	}

	t := NewTable("ID", "Name", "GCID", "Virtual Balance", "Mirror ID", "Created", "Tokens")
	for _, p := range portfolios {
		t.AppendRow(table.Row{
			truncateID(p.ID),
			Cyan(p.Name),
			p.GCID,
			FormatMoney(p.VirtualBalance),
			p.MirrorID,
			formatTimestamp(p.CreatedAt),
			p.TokenCount,
		})
	}
	RenderTable(t)
}

type AgentPortfolioDetail struct {
	ID             string
	Name           string
	GCID           int
	VirtualBalance float64
	MirrorID       int
	CreatedAt      string
	Tokens         []AgentPortfolioTokenRow
}

type AgentPortfolioTokenRow struct {
	TokenID   string
	Name      string
	ClientID  string
	AppName   string
	Scopes    string
	IPs       string
	ExpiresAt string
	CreatedAt string
}

func PrintAgentPortfolioDetail(detail AgentPortfolioDetail, format Format) {
	if format == JSON {
		PrintJSON(detail)
		return
	}

	t := NewDetailTable()
	DetailRow(t, "Portfolio ID", detail.ID)
	DetailRow(t, "Name", Cyan(detail.Name))
	DetailRow(t, "GCID", detail.GCID)
	DetailRow(t, "Virtual Balance", FormatMoney(detail.VirtualBalance))
	DetailRow(t, "Mirror ID", detail.MirrorID)
	if detail.CreatedAt != "" {
		DetailRow(t, "Created", formatTimestamp(detail.CreatedAt))
	}
	RenderTable(t)

	if len(detail.Tokens) > 0 {
		fmt.Printf("\n  %s\n\n", Bold("User Tokens"))
		tt := NewTable("Token ID", "Name", "Scopes", "IPs", "Expires", "Created")
		for _, tok := range detail.Tokens {
			tt.AppendRow(table.Row{
				truncateID(tok.TokenID),
				tok.Name,
				tok.Scopes,
				tok.IPs,
				formatTimestamp(tok.ExpiresAt),
				formatTimestamp(tok.CreatedAt),
			})
		}
		RenderTable(tt)
	}
}

func PrintCreatedAgentPortfolio(id, name string, gcid int, virtualBalance float64, mirrorID int, tokens []CreatedTokenDisplay, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"status":         "created",
			"portfolioId":    id,
			"name":           name,
			"gcid":           gcid,
			"virtualBalance": virtualBalance,
			"mirrorId":       mirrorID,
			"tokens":         tokens,
		})
		return
	}

	Successf("Agent portfolio %q created successfully.", name)
	fmt.Println()

	t := NewDetailTable()
	DetailRow(t, "Portfolio ID", id)
	DetailRow(t, "Name", Cyan(name))
	DetailRow(t, "GCID", gcid)
	DetailRow(t, "Virtual Balance", FormatMoney(virtualBalance))
	DetailRow(t, "Mirror ID", mirrorID)
	RenderTable(t)

	for _, tok := range tokens {
		fmt.Printf("\n  %s\n\n", Bold("User Token"))
		tt := NewDetailTable()
		DetailRow(tt, "Token ID", tok.TokenID)
		DetailRow(tt, "Token Name", tok.TokenName)
		DetailRow(tt, "Token Secret", Red(tok.TokenSecret))
		DetailRow(tt, "Scopes", tok.Scopes)
		RenderTable(tt)
		Warnf("Save the token secret now — it will not be shown again.")
		fmt.Printf("\n  %s\n", Bold("Quick start — trade on this agent portfolio:"))
		fmt.Printf("  ETORO_USER_KEY=%s etoro trade open BTC --amount 50 --yes\n", tok.TokenSecret)
	}
}

type CreatedTokenDisplay struct {
	TokenID     string `json:"tokenId"`
	TokenName   string `json:"tokenName"`
	TokenSecret string `json:"tokenSecret"`
	Scopes      string `json:"scopes"`
}

func PrintCreatedUserToken(tokenID, tokenSecret string, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"status":      "created",
			"userTokenId": tokenID,
			"userToken":   tokenSecret,
		})
		return
	}

	Successf("User token created successfully.")
	fmt.Println()
	t := NewDetailTable()
	DetailRow(t, "Token ID", tokenID)
	DetailRow(t, "Token Secret", Red(tokenSecret))
	RenderTable(t)
	Warnf("Save the token secret now — it will not be shown again.")
}

func FormatScopeIDs(ids []int) string {
	scopeNames := map[int]string{
		200: "real:read",
		201: "demo:read",
		202: "real:write",
		203: "demo:write",
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		if name, ok := scopeNames[id]; ok {
			parts = append(parts, name)
		} else {
			parts = append(parts, fmt.Sprintf("%d", id))
		}
	}
	return strings.Join(parts, ", ")
}

func truncateID(id string) string {
	if len(id) > 12 {
		return id[:8] + "..."
	}
	return id
}

func formatTimestamp(ts string) string {
	if ts == "" {
		return "—"
	}
	if len(ts) > 19 {
		return ts[:19]
	}
	return ts
}
