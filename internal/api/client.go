package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/marianopa-tr/etoro-cli/internal/config"
	"github.com/google/uuid"
)

const UserAgent = "etoro-cli/1.0"

type Client struct {
	cfg        *config.Config
	httpClient *http.Client
	demo       bool
	baseURL    string
}

func NewClient(cfg *config.Config, demo bool) *Client {
	timeout := cfg.TimeoutDuration()
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: timeout},
		demo:       demo,
		baseURL:    config.DefaultBaseURL,
	}
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) SetTimeout(d time.Duration) {
	c.httpClient.Timeout = d
}

func (c *Client) IsDemo() bool {
	return c.demo
}

func (c *Client) doRequest(method, path string, body any) (*http.Response, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("x-request-id", uuid.New().String())
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.cfg.Auth.APIKey != "" {
		req.Header.Set("x-api-key", c.cfg.Auth.APIKey)
	}
	if c.cfg.Auth.UserKey != "" {
		req.Header.Set("x-user-key", c.cfg.Auth.UserKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}

func (c *Client) get(path string, params map[string]string) ([]byte, error) {
	url := c.baseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("x-request-id", uuid.New().String())
	if c.cfg.Auth.APIKey != "" {
		req.Header.Set("x-api-key", c.cfg.Auth.APIKey)
	}
	if c.cfg.Auth.UserKey != "" {
		req.Header.Set("x-user-key", c.cfg.Auth.UserKey)
	}

	if len(params) > 0 {
		parts := make([]string, 0, len(params))
		for k, v := range params {
			parts = append(parts, neturl.QueryEscape(k)+"="+neturl.QueryEscape(v))
		}
		raw := strings.Join(parts, "&")
		raw = strings.ReplaceAll(raw, "%2C", ",")
		req.URL.RawQuery = raw
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(data),
		}
	}

	return data, nil
}

func (c *Client) post(path string, body any) ([]byte, error) {
	resp, err := c.doRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	return data, nil
}

func (c *Client) put(path string, body any) ([]byte, error) {
	resp, err := c.doRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	return data, nil
}

func (c *Client) patch(path string, body any) ([]byte, error) {
	resp, err := c.doRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	return data, nil
}

func (c *Client) delete(path string) ([]byte, error) {
	resp, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(data)}
	}

	return data, nil
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("API error (HTTP %d)", e.StatusCode)
}
