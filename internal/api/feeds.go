package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (c *Client) GetInstrumentFeed(instrumentID int, offset, limit int) (*FeedResponse, error) {
	params := map[string]string{}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		params["take"] = strconv.Itoa(limit)
	}

	data, err := c.get(fmt.Sprintf("/api/v1/feeds/instrument/%d", instrumentID), params)
	if err != nil {
		return nil, err
	}

	var resp FeedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing instrument feed response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetUserFeed(userID string, offset, limit int) (*FeedResponse, error) {
	params := map[string]string{}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		params["take"] = strconv.Itoa(limit)
	}

	data, err := c.get(fmt.Sprintf("/api/v1/feeds/user/%s", userID), params)
	if err != nil {
		return nil, err
	}

	var resp FeedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing user feed response: %w", err)
	}
	return &resp, nil
}

func (c *Client) CreatePost(req *CreatePostRequest) (json.RawMessage, error) {
	data, err := c.post("/api/v1/feeds/post", req)
	if err != nil {
		return nil, err
	}
	return data, nil
}
