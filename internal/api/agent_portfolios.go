package api

import (
	"encoding/json"
	"fmt"
)

const agentPortfoliosPath = "/api/v1/agent-portfolios"

func (c *Client) CreateAgentPortfolio(req *CreateAgentPortfolioRequest) (*CreateAgentPortfolioResponse, error) {
	data, err := c.post(agentPortfoliosPath, req)
	if err != nil {
		return nil, err
	}

	var resp CreateAgentPortfolioResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing create agent-portfolio response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetAgentPortfolios() (*GetAgentPortfoliosResponse, error) {
	data, err := c.get(agentPortfoliosPath, nil)
	if err != nil {
		return nil, err
	}

	var resp GetAgentPortfoliosResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing agent-portfolios response: %w", err)
	}
	return &resp, nil
}

func (c *Client) DeleteAgentPortfolio(portfolioID string) error {
	_, err := c.delete(fmt.Sprintf("%s/%s", agentPortfoliosPath, portfolioID))
	return err
}

func (c *Client) CreateUserToken(portfolioID string, req *CreateUserTokenRequest) (*CreateUserTokenResponse, error) {
	path := fmt.Sprintf("%s/%s/user-tokens", agentPortfoliosPath, portfolioID)
	data, err := c.post(path, req)
	if err != nil {
		return nil, err
	}

	var resp CreateUserTokenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing create user-token response: %w", err)
	}
	return &resp, nil
}

func (c *Client) UpdateUserToken(portfolioID, tokenID string, req *UpdateUserTokenRequest) error {
	path := fmt.Sprintf("%s/%s/user-tokens/%s", agentPortfoliosPath, portfolioID, tokenID)
	_, err := c.patch(path, req)
	return err
}

func (c *Client) DeleteUserToken(portfolioID, tokenID string) error {
	path := fmt.Sprintf("%s/%s/user-tokens/%s", agentPortfoliosPath, portfolioID, tokenID)
	_, err := c.delete(path)
	return err
}
