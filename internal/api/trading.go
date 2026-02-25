package api

import (
	"encoding/json"
	"fmt"
)

func (c *Client) tradingPrefix() string {
	if c.demo {
		return "/api/v1/trading/execution/demo"
	}
	return "/api/v1/trading/execution"
}

func (c *Client) OpenPositionByAmount(req *OpenByAmountRequest) (json.RawMessage, error) {
	path := c.tradingPrefix() + "/market-open-orders/by-amount"
	data, err := c.post(path, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) OpenPositionByUnits(req *OpenByUnitsRequest) (json.RawMessage, error) {
	path := c.tradingPrefix() + "/market-open-orders/by-units"
	data, err := c.post(path, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) PlaceLimitOrder(req *LimitOrderRequest) (json.RawMessage, error) {
	path := c.tradingPrefix() + "/limit-orders"
	data, err := c.post(path, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) ClosePosition(positionID int, req *ClosePositionRequest) (json.RawMessage, error) {
	path := fmt.Sprintf("%s/market-close-orders/positions/%d", c.tradingPrefix(), positionID)
	data, err := c.post(path, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) CancelOrder(orderID int) (json.RawMessage, error) {
	path := fmt.Sprintf("%s/market-open-orders/%d", c.tradingPrefix(), orderID)
	data, err := c.delete(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) CancelLimitOrder(orderID int) (json.RawMessage, error) {
	path := fmt.Sprintf("%s/limit-orders/%d", c.tradingPrefix(), orderID)
	data, err := c.delete(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
