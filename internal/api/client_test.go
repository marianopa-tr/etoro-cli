package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/marianopa-tr/etoro-cli/internal/config"
)

func newTestClient(handler http.HandlerFunc) (*httptest.Server, *Client) {
	srv := httptest.NewServer(handler)
	cfg := &config.Config{
		Auth: config.AuthConfig{
			APIKey:  "test-api-key",
			UserKey: "test-user-key",
		},
		Defaults: config.DefaultsConfig{Timeout: "10s"},
	}
	client := NewClient(cfg, false)
	client.SetBaseURL(srv.URL)
	return srv, client
}

func TestClientSetsAuthHeaders(t *testing.T) {
	var gotHeaders http.Header
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.get("/api/v1/test", nil)
	if err != nil {
		t.Fatalf("get() error: %v", err)
	}

	if got := gotHeaders.Get("x-api-key"); got != "test-api-key" {
		t.Errorf("x-api-key = %q, want %q", got, "test-api-key")
	}
	if got := gotHeaders.Get("x-user-key"); got != "test-user-key" {
		t.Errorf("x-user-key = %q, want %q", got, "test-user-key")
	}
	if got := gotHeaders.Get("x-request-id"); got == "" {
		t.Error("x-request-id should be set")
	}
	if got := gotHeaders.Get("User-Agent"); got != UserAgent {
		t.Errorf("User-Agent = %q, want %q", got, UserAgent)
	}
}

func TestClientOmitsEmptyAuthHeaders(t *testing.T) {
	var gotHeaders http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cfg := &config.Config{
		Auth:     config.AuthConfig{},
		Defaults: config.DefaultsConfig{Timeout: "5s"},
	}
	client := NewClient(cfg, false)
	client.SetBaseURL(srv.URL)

	_, err := client.get("/test", nil)
	if err != nil {
		t.Fatalf("get() error: %v", err)
	}

	if gotHeaders.Get("x-api-key") != "" {
		t.Error("x-api-key should not be set when empty")
	}
	if gotHeaders.Get("x-user-key") != "" {
		t.Error("x-user-key should not be set when empty")
	}
}

func TestClientGetQueryParams(t *testing.T) {
	var gotQuery string
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.get("/test", map[string]string{"foo": "bar", "baz": "qux"})
	if err != nil {
		t.Fatalf("get() error: %v", err)
	}

	if !strings.Contains(gotQuery, "foo=bar") {
		t.Errorf("query should contain 'foo=bar', got %q", gotQuery)
	}
	if !strings.Contains(gotQuery, "baz=qux") {
		t.Errorf("query should contain 'baz=qux', got %q", gotQuery)
	}
}

func TestClientPostSendsBody(t *testing.T) {
	var gotBody map[string]any
	var gotMethod, gotContentType string
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Write([]byte(`{"orderId":123}`))
	})
	defer srv.Close()

	_, err := client.post("/test", map[string]any{"InstrumentID": 1001, "Amount": 500.0})
	if err != nil {
		t.Fatalf("post() error: %v", err)
	}

	if gotMethod != "POST" {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", gotContentType)
	}
	if gotBody["InstrumentID"] != float64(1001) {
		t.Errorf("InstrumentID = %v, want 1001", gotBody["InstrumentID"])
	}
}

func TestClientPut(t *testing.T) {
	var gotMethod string
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.put("/test", map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("put() error: %v", err)
	}
	if gotMethod != "PUT" {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
}

func TestClientDelete(t *testing.T) {
	var gotMethod string
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.delete("/test")
	if err != nil {
		t.Fatalf("delete() error: %v", err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

func TestClientHandles4xxError(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"errorCode":"Unauthorized"}`))
	})
	defer srv.Close()

	_, err := client.get("/test", nil)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Body, "Unauthorized") {
		t.Errorf("Body should contain 'Unauthorized', got %q", apiErr.Body)
	}
}

func TestClientHandles5xxError(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`Internal Server Error`))
	})
	defer srv.Close()

	_, err := client.post("/test", nil)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestClientIsDemo(t *testing.T) {
	cfg := &config.Config{Defaults: config.DefaultsConfig{Timeout: "5s"}}

	if NewClient(cfg, false).IsDemo() {
		t.Error("IsDemo() should be false for real client")
	}
	if !NewClient(cfg, true).IsDemo() {
		t.Error("IsDemo() should be true for demo client")
	}
}

func TestClientSetTimeout(t *testing.T) {
	cfg := &config.Config{Defaults: config.DefaultsConfig{Timeout: "5s"}}
	client := NewClient(cfg, false)
	client.SetTimeout(10 * time.Second)
	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", client.httpClient.Timeout)
	}
}

func TestSetBaseURL(t *testing.T) {
	cfg := &config.Config{Defaults: config.DefaultsConfig{Timeout: "5s"}}
	client := NewClient(cfg, false)
	if client.baseURL != config.DefaultBaseURL {
		t.Errorf("default baseURL = %q, want %q", client.baseURL, config.DefaultBaseURL)
	}
	client.SetBaseURL("http://localhost:9999")
	if client.baseURL != "http://localhost:9999" {
		t.Errorf("baseURL = %q after SetBaseURL", client.baseURL)
	}
}

func TestAPIErrorString(t *testing.T) {
	err := &APIError{StatusCode: 404, Body: "not found"}
	if !strings.Contains(err.Error(), "404") || !strings.Contains(err.Error(), "not found") {
		t.Errorf("Error() = %q", err.Error())
	}

	errNoBody := &APIError{StatusCode: 500}
	if !strings.Contains(errNoBody.Error(), "500") {
		t.Errorf("Error() = %q", errNoBody.Error())
	}
}

func TestSearchInstruments(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("internalSymbolFull") != "AAPL" {
			t.Errorf("internalSymbolFull = %q", r.URL.Query().Get("internalSymbolFull"))
		}
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields should be set")
		}
		json.NewEncoder(w).Encode(InstrumentSearchResponse{
			TotalItems: 1,
			Items:      []Instrument{{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"}},
		})
	})
	defer srv.Close()

	resp, err := client.SearchInstruments("AAPL", 1, 10)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Symbol != "AAPL" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetRates(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LiveRatesResponse{
			Rates: []Rate{{InstrumentID: 1001, Ask: 150.25, Bid: 150.20}},
		})
	})
	defer srv.Close()

	resp, err := client.GetRates([]int{1001})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Rates) != 1 || resp.Rates[0].Ask != 150.25 {
		t.Errorf("unexpected: %+v", resp)
	}
}

func TestGetInstruments(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(InstrumentsResponse{
			InstrumentDisplayDatas: []InstrumentDisplayData{
				{InstrumentID: 1001, Symbol: "AAPL"},
				{InstrumentID: 1002, Symbol: "TSLA"},
			},
		})
	})
	defer srv.Close()

	resp, err := client.GetInstruments([]int{1001, 1002})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.InstrumentDisplayDatas) != 2 {
		t.Errorf("count = %d", len(resp.InstrumentDisplayDatas))
	}
}

func TestGetExchanges(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ExchangesResponse{
			ExchangeInfo: []Exchange{{ExchangeID: 1, Name: "NASDAQ"}},
		})
	})
	defer srv.Close()

	resp, err := client.GetExchanges()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.ExchangeInfo[0].Name != "NASDAQ" {
		t.Errorf("Name = %q", resp.ExchangeInfo[0].Name)
	}
}

func TestGetInstrumentTypes(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(InstrumentTypesResponse{
			InstrumentTypes: []InstrumentTypeInfo{{InstrumentTypeID: 1, Name: "Stocks"}},
		})
	})
	defer srv.Close()

	resp, err := client.GetInstrumentTypes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.InstrumentTypes[0].Name != "Stocks" {
		t.Errorf("Name = %q", resp.InstrumentTypes[0].Name)
	}
}

func TestTradingPrefixes(t *testing.T) {
	cfg := &config.Config{Defaults: config.DefaultsConfig{Timeout: "5s"}}

	real := NewClient(cfg, false)
	if real.tradingPrefix() != "/api/v1/trading/execution" {
		t.Errorf("real trading = %q", real.tradingPrefix())
	}
	if real.infoPrefix() != "/api/v1/trading/info" {
		t.Errorf("real info = %q", real.infoPrefix())
	}

	demo := NewClient(cfg, true)
	if demo.tradingPrefix() != "/api/v1/trading/execution/demo" {
		t.Errorf("demo trading = %q", demo.tradingPrefix())
	}
	if demo.infoPrefix() != "/api/v1/trading/info/demo" {
		t.Errorf("demo info = %q", demo.infoPrefix())
	}
}

func TestGetPortfolio(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/info/portfolio" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(PortfolioResponse{
			ClientPortfolio: ClientPortfolio{
				Credit:    10000,
				Positions: []Position{{PositionID: 1, InstrumentID: 1001, Amount: 500}},
			},
		})
	})
	defer srv.Close()

	resp, err := client.GetPortfolio()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.ClientPortfolio.Credit != 10000 {
		t.Errorf("Credit = %f", resp.ClientPortfolio.Credit)
	}
	if len(resp.ClientPortfolio.Positions) != 1 {
		t.Errorf("Positions = %d", len(resp.ClientPortfolio.Positions))
	}
}

func TestGetPnL(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/info/real/pnl" {
			t.Errorf("path = %q, want /api/v1/trading/info/real/pnl", r.URL.Path)
		}
		json.NewEncoder(w).Encode(PortfolioResponse{
			ClientPortfolio: ClientPortfolio{UnrealizedPnL: 250.50},
		})
	})
	defer srv.Close()

	resp, err := client.GetPnL()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.ClientPortfolio.UnrealizedPnL != 250.50 {
		t.Errorf("PnL = %f", resp.ClientPortfolio.UnrealizedPnL)
	}
}

func TestOpenPositionByAmount(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/trading/execution/market-open-orders/by-amount" {
			t.Errorf("method=%s path=%s", r.Method, r.URL.Path)
		}
		var body OpenByAmountRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.InstrumentID != 1001 || body.Amount != 500 || !body.IsBuy {
			t.Errorf("body = %+v", body)
		}
		w.Write([]byte(`{"orderId":999}`))
	})
	defer srv.Close()

	_, err := client.OpenPositionByAmount(&OpenByAmountRequest{
		InstrumentID: 1001, IsBuy: true, Leverage: 1, Amount: 500,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestOpenPositionByUnits(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/execution/market-open-orders/by-units" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Write([]byte(`{"orderId":998}`))
	})
	defer srv.Close()

	_, err := client.OpenPositionByUnits(&OpenByUnitsRequest{
		InstrumentID: 1001, IsBuy: true, Leverage: 1, Units: 0.5,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestPlaceLimitOrder(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/execution/limit-orders" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Write([]byte(`{"orderId":997}`))
	})
	defer srv.Close()

	_, err := client.PlaceLimitOrder(&LimitOrderRequest{
		InstrumentID: 1001, IsBuy: true, Leverage: 1, Amount: 500, Rate: 150,
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestClosePosition(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/execution/market-close-orders/positions/12345" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Write([]byte(`{"status":"closed"}`))
	})
	defer srv.Close()

	_, err := client.ClosePosition(12345, &ClosePositionRequest{})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestCancelOrder(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s", r.Method)
		}
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.CancelOrder(555)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestCancelLimitOrder(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" || !strings.Contains(r.URL.Path, "limit-orders/666") {
			t.Errorf("method=%s path=%s", r.Method, r.URL.Path)
		}
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	_, err := client.CancelLimitOrder(666)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestGetWatchlists(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(WatchlistsResponse{
			Watchlists: []Watchlist{{WatchlistID: "abc", Name: "Tech", TotalItems: 5}},
		})
	})
	defer srv.Close()

	resp, err := client.GetWatchlists()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Watchlists) != 1 || resp.Watchlists[0].Name != "Tech" {
		t.Errorf("unexpected: %+v", resp)
	}
}

func TestCreateWatchlist(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/v1/watchlists" {
			t.Errorf("method=%s path=%s", r.Method, r.URL.Path)
		}
		w.Write([]byte(`{"watchlistId":"new-id"}`))
	})
	defer srv.Close()

	_, err := client.CreateWatchlist(&CreateWatchlistRequest{Name: "Test"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestGetTradeHistory(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trading/info/trade/history" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("minDate") != "2025-01-01" {
			t.Errorf("minDate = %q", r.URL.Query().Get("minDate"))
		}
		json.NewEncoder(w).Encode([]TradeHistoryEntry{
			{PositionID: 100, InstrumentID: 1001, NetProfit: 50.0},
		})
	})
	defer srv.Close()

	entries, err := client.GetTradeHistory("2025-01-01", 1, 50)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(entries) != 1 || entries[0].NetProfit != 50.0 {
		t.Errorf("unexpected: %+v", entries)
	}
}

func TestGetUserGain(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/people/trader1/gain") {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(UserGainResponse{
			Monthly: []GainEntry{{Timestamp: "2025-01", Gain: 5.5}},
			Yearly:  []GainEntry{{Timestamp: "2024", Gain: 25.0}},
		})
	})
	defer srv.Close()

	resp, err := client.GetUserGain("trader1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Monthly) != 1 || resp.Monthly[0].Gain != 5.5 {
		t.Errorf("unexpected monthly: %+v", resp.Monthly)
	}
}

func TestGetCopiers(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(CopiersResponse{
			Copiers: []CopierInfo{{Country: "US", Club: "Gold"}},
		})
	})
	defer srv.Close()

	resp, err := client.GetCopiers()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Copiers) != 1 || resp.Copiers[0].Country != "US" {
		t.Errorf("unexpected: %+v", resp)
	}
}

func TestGetInstrumentFeed(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/feeds/instrument/1001") {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(FeedResponse{
			Discussions: []FeedPost{{ID: "p1", Post: PostData{Message: PostMessage{Text: "Bullish!"}}}},
		})
	})
	defer srv.Close()

	resp, err := client.GetInstrumentFeed(1001, 0, 10)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.Discussions) != 1 {
		t.Errorf("Discussions = %d", len(resp.Discussions))
	}
}

func TestGetUserFeed(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/feeds/user/42") {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(FeedResponse{Discussions: []FeedPost{}})
	})
	defer srv.Close()

	_, err := client.GetUserFeed("42", 0, 10)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestSearchUsers(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("period") != "LastYear" {
			t.Errorf("period = %q", r.URL.Query().Get("period"))
		}
		w.Write([]byte(`{"items":[]}`))
	})
	defer srv.Close()

	_, err := client.SearchUsers("LastYear", 1, 20)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestGetCuratedLists(t *testing.T) {
	srv, client := newTestClient(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(CuratedListsResponse{
			CuratedLists: []CuratedList{{Name: "Top Tech", Description: "Top tech stocks"}},
		})
	})
	defer srv.Close()

	resp, err := client.GetCuratedLists()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp.CuratedLists) != 1 {
		t.Errorf("count = %d", len(resp.CuratedLists))
	}
}
