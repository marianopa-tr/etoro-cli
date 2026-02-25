package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetWatchlists() (*WatchlistsResponse, error) {
	data, err := c.get("/api/v1/watchlists", nil)
	if err != nil {
		return nil, err
	}

	var resp WatchlistsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing watchlists response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetWatchlist(id string) (*Watchlist, error) {
	data, err := c.get("/api/v1/watchlists/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Watchlist
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing watchlist response: %w", err)
	}
	return &resp, nil
}

func (c *Client) CreateWatchlist(req *CreateWatchlistRequest) (json.RawMessage, error) {
	data, err := c.post("/api/v1/watchlists", req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) DeleteWatchlist(id string) error {
	_, err := c.delete("/api/v1/watchlists/" + id)
	return err
}

func (c *Client) AddWatchlistItems(watchlistID string, items []WatchlistItem) (json.RawMessage, error) {
	data, err := c.post(fmt.Sprintf("/api/v1/watchlists/%s/items", watchlistID), items)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) RemoveWatchlistItems(watchlistID string, items []WatchlistItem) error {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/watchlists/%s/items", watchlistID), items)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to remove items (HTTP %d)", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetCuratedLists() (*CuratedListsResponse, error) {
	data, err := c.get("/api/v1/curated-lists", nil)
	if err != nil {
		return nil, err
	}

	var resp CuratedListsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing curated lists response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetDefaultWatchlistItems() (json.RawMessage, error) {
	data, err := c.get("/api/v1/watchlists/default-watchlists/items", nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}
