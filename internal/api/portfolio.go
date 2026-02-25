package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (c *Client) infoPrefix() string {
	if c.demo {
		return "/api/v1/trading/info/demo"
	}
	return "/api/v1/trading/info"
}

func (c *Client) GetPortfolio() (*PortfolioResponse, error) {
	var path string
	if c.demo {
		path = "/api/v1/trading/info/demo/portfolio"
	} else {
		path = "/api/v1/trading/info/portfolio"
	}

	data, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}

	var resp PortfolioResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing portfolio response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetPnL() (*PortfolioResponse, error) {
	var path string
	if c.demo {
		path = "/api/v1/trading/info/demo/pnl"
	} else {
		path = "/api/v1/trading/info/real/pnl"
	}

	data, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}

	var resp PortfolioResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing PnL response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetTradeHistory(minDate string, page, pageSize int) ([]TradeHistoryEntry, error) {
	if c.demo {
		return nil, fmt.Errorf("trade history is not available for demo accounts (eToro API limitation). Use a real-account key or view demo positions with: etoro portfolio positions --demo")
	}

	params := map[string]string{
		"minDate": minDate,
	}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}

	data, err := c.get("/api/v1/trading/info/trade/history", params)
	if err != nil {
		return nil, err
	}

	var entries []TradeHistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing trade history response: %w", err)
	}
	return entries, nil
}
