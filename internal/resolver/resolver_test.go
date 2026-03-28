package resolver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/config"
)

func newTestResolver(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Resolver) {
	t.Helper()
	srv := httptest.NewServer(handler)
	cfg := &config.Config{
		Auth:     config.AuthConfig{APIKey: "k", UserKey: "u"},
		Defaults: config.DefaultsConfig{Timeout: "5s"},
	}
	client := api.NewClient(cfg, false)
	client.SetBaseURL(srv.URL)

	t.Setenv("ETORO_CACHE_DIR", t.TempDir())

	r := New(client)
	return srv, r
}

func TestResolveByID(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		t.Error("should not call API for numeric ID")
	})
	defer srv.Close()

	id, symbol, err := r.Resolve("1001")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if id != 1001 {
		t.Errorf("id = %d, want 1001", id)
	}
	if symbol != "1001" {
		t.Errorf("symbol = %q, want '1001'", symbol)
	}
}

func TestResolveBySymbol(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"},
			},
		})
	})
	defer srv.Close()

	id, symbol, err := r.Resolve("AAPL")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if id != 1001 {
		t.Errorf("id = %d, want 1001", id)
	}
	if symbol != "AAPL" {
		t.Errorf("symbol = %q, want AAPL", symbol)
	}
}

func TestResolveExactMatch(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 2002, Symbol: "AAPLX", DisplayName: "Apple Extended"},
				{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"},
			},
		})
	})
	defer srv.Close()

	id, symbol, err := r.Resolve("AAPL")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if id != 1001 || symbol != "AAPL" {
		t.Errorf("id=%d symbol=%q, want 1001/AAPL (exact match preferred)", id, symbol)
	}
}

func TestResolveCaseInsensitive(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"},
			},
		})
	})
	defer srv.Close()

	id, _, err := r.Resolve("aapl")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if id != 1001 {
		t.Errorf("id = %d, want 1001", id)
	}
}

func TestResolveNotFound(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{Items: []api.Instrument{}})
	})
	defer srv.Close()

	_, _, err := r.Resolve("NOTEXIST")
	if err == nil {
		t.Fatal("expected error for non-existent instrument")
	}
}

func TestResolveFallsBackToFirstResult(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 3003, Symbol: "BTC", DisplayName: "Bitcoin"},
			},
		})
	})
	defer srv.Close()

	id, symbol, err := r.Resolve("bitcoin")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if id != 3003 || symbol != "BTC" {
		t.Errorf("id=%d symbol=%q, want 3003/BTC", id, symbol)
	}
}

func TestResolveUsesCache(t *testing.T) {
	callCount := 0
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		callCount++
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"},
			},
		})
	})
	defer srv.Close()

	r.Resolve("AAPL")
	r.Resolve("AAPL")

	if callCount != 1 {
		t.Errorf("API called %d times, want 1 (should use cache)", callCount)
	}
}

func TestResolveMultiple(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		sym := req.URL.Query().Get("internalSymbolFull")
		var items []api.Instrument
		switch sym {
		case "AAPL":
			items = []api.Instrument{{InstrumentID: 1001, Symbol: "AAPL"}}
		case "TSLA":
			items = []api.Instrument{{InstrumentID: 1002, Symbol: "TSLA"}}
		}
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{Items: items})
	})
	defer srv.Close()

	ids, symbols, err := r.ResolveMultiple([]string{"AAPL", "TSLA"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("ids count = %d", len(ids))
	}
	if ids[0] != 1001 || ids[1] != 1002 {
		t.Errorf("ids = %v", ids)
	}
	if symbols[0] != "AAPL" || symbols[1] != "TSLA" {
		t.Errorf("symbols = %v", symbols)
	}
}

func TestCachedSymbols(t *testing.T) {
	srv, r := newTestResolver(t, func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(api.InstrumentSearchResponse{
			Items: []api.Instrument{
				{InstrumentID: 1001, Symbol: "AAPL", DisplayName: "Apple Inc"},
			},
		})
	})
	defer srv.Close()

	r.Resolve("AAPL")
	symbols := r.CachedSymbols()
	if len(symbols) != 1 || symbols[0] != "AAPL" {
		t.Errorf("CachedSymbols = %v", symbols)
	}
}

func TestCachePersistence(t *testing.T) {
	cacheDir := t.TempDir()

	cacheFile := filepath.Join(cacheDir, "instruments.json")
	cacheData := map[string]cachedEntry{
		"AAPL": {
			InstrumentID: 1001,
			Symbol:       "AAPL",
			DisplayName:  "Apple Inc",
			CachedAt:     9999999999,
		},
	}
	data, _ := json.Marshal(cacheData)
	os.MkdirAll(cacheDir, 0o700)
	os.WriteFile(cacheFile, data, 0o600)

	origCacheDir := os.Getenv("HOME")
	r := &Resolver{
		cache: make(map[string]cachedEntry),
	}

	_ = origCacheDir
	r.cache = cacheData
	r.dirty = true
	r.saveCache()

	if len(r.CachedSymbols()) != 1 {
		t.Errorf("expected 1 cached symbol, got %d", len(r.CachedSymbols()))
	}
}
