package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const defaultSearchFields = "instrumentId,displayname,internalSymbolFull,instrumentType,exchangeID,isExchangeOpen,isCurrentlyTradable,isBuyEnabled,currentRate,dailyPriceChange,weeklyPriceChange,monthlyPriceChange,threeMonthPriceChange,oneYearPriceChange,popularityUniques7Day,internalAssetClassName,internalExchangeName,isHiddenFromClient"

func (c *Client) SearchInstruments(query string, page, pageSize int) (*InstrumentSearchResponse, error) {
	params := map[string]string{
		"fields": defaultSearchFields,
	}
	if query != "" {
		params["internalSymbolFull"] = query
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	if page > 0 {
		params["pageNumber"] = strconv.Itoa(page)
	}

	data, err := c.get("/api/v1/market-data/search", params)
	if err != nil {
		return nil, err
	}

	var resp InstrumentSearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}
	return &resp, nil
}

func (c *Client) SearchInstrumentsByText(query string, page, pageSize int) (*InstrumentSearchResponse, error) {
	params := map[string]string{
		"fields": defaultSearchFields,
	}
	if query != "" {
		params["displayname"] = query
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	if page > 0 {
		params["pageNumber"] = strconv.Itoa(page)
	}

	data, err := c.get("/api/v1/market-data/search", params)
	if err != nil {
		return nil, err
	}

	var resp InstrumentSearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing search response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetInstruments(ids []int) (*InstrumentsResponse, error) {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = strconv.Itoa(id)
	}
	params := map[string]string{
		"instrumentIds": strings.Join(strs, ","),
	}

	data, err := c.get("/api/v1/market-data/instruments", params)
	if err != nil {
		return nil, err
	}

	var resp InstrumentsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing instruments response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetRates(ids []int) (*LiveRatesResponse, error) {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = strconv.Itoa(id)
	}
	params := map[string]string{
		"instrumentIds": strings.Join(strs, ","),
	}

	data, err := c.get("/api/v1/market-data/instruments/rates", params)
	if err != nil {
		return nil, err
	}

	var resp LiveRatesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing rates response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetExchanges() (*ExchangesResponse, error) {
	data, err := c.get("/api/v1/market-data/exchanges", nil)
	if err != nil {
		return nil, err
	}

	var resp ExchangesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing exchanges response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetInstrumentTypes() (*InstrumentTypesResponse, error) {
	data, err := c.get("/api/v1/market-data/instrument-types", nil)
	if err != nil {
		return nil, err
	}

	var resp InstrumentTypesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing instrument types response: %w", err)
	}
	return &resp, nil
}
