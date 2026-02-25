package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) SearchUsers(period string, page, pageSize int) (json.RawMessage, error) {
	params := map[string]string{
		"period": period,
	}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if pageSize > 0 {
		params["pageSize"] = fmt.Sprintf("%d", pageSize)
	}

	data, err := c.get("/api/v1/user-info/people/search", params)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetUserGain(username string) (*UserGainResponse, error) {
	data, err := c.get(fmt.Sprintf("/api/v1/user-info/people/%s/gain", username), nil)
	if err != nil {
		return nil, err
	}

	var resp UserGainResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing user gain response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetUserDailyGain(username string) (json.RawMessage, error) {
	data, err := c.get(fmt.Sprintf("/api/v1/user-info/people/%s/daily-gain", username), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetUserPortfolioLive(username string) (json.RawMessage, error) {
	data, err := c.get(fmt.Sprintf("/api/v1/user-info/people/%s/portfolio/live", username), nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetUserTradeInfo(username, period string) (*TradeInfoResponse, error) {
	params := map[string]string{
		"period": period,
	}
	data, err := c.get(fmt.Sprintf("/api/v1/user-info/people/%s/tradeinfo", username), params)
	if err != nil {
		return nil, err
	}
	var resp TradeInfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing tradeinfo response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetPeople() (json.RawMessage, error) {
	data, err := c.get("/api/v1/user-info/people", nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetUserProfile(username string) (*UserProfile, error) {
	params := map[string]string{
		"usernames": username,
	}
	data, err := c.get("/api/v1/user-info/people", params)
	if err != nil {
		return nil, err
	}
	var resp UserProfileResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing user profile response: %w", err)
	}
	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("user %q not found", username)
	}
	return &resp.Users[0], nil
}

func (c *Client) GetCopiers() (*CopiersResponse, error) {
	data, err := c.get("/api/v1/pi-data/copiers", nil)
	if err != nil {
		return nil, err
	}

	var resp CopiersResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing copiers response: %w", err)
	}
	return &resp, nil
}
