package resolver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/marianopa-tr/etoro-cli/internal/api"
	"github.com/marianopa-tr/etoro-cli/internal/config"
)

type Resolver struct {
	client *api.Client
	cache  map[string]cachedEntry
	dirty  bool
}

type cachedEntry struct {
	InstrumentID int    `json:"instrumentId"`
	Symbol       string `json:"symbol"`
	DisplayName  string `json:"displayName"`
	CachedAt     int64  `json:"cachedAt"`
}

const cacheTTL = 24 * time.Hour

func New(client *api.Client) *Resolver {
	r := &Resolver{
		client: client,
		cache:  make(map[string]cachedEntry),
	}
	r.loadCache()
	return r
}

func (r *Resolver) Resolve(symbolOrID string) (int, string, error) {
	if id, err := strconv.Atoi(symbolOrID); err == nil {
		return id, symbolOrID, nil
	}

	key := strings.ToUpper(symbolOrID)
	if entry, ok := r.cache[key]; ok {
		if time.Since(time.Unix(entry.CachedAt, 0)) < cacheTTL {
			return entry.InstrumentID, entry.Symbol, nil
		}
	}

	resp, err := r.client.SearchInstruments(symbolOrID, 1, 50)
	if err != nil {
		return 0, "", fmt.Errorf("searching for %q: %w", symbolOrID, err)
	}

	if len(resp.Items) == 0 {
		return 0, "", fmt.Errorf("no instrument found for %q", symbolOrID)
	}

	for _, inst := range resp.Items {
		if strings.EqualFold(inst.Symbol, symbolOrID) {
			r.cacheInstrument(inst)
			return inst.InstrumentID, inst.Symbol, nil
		}
	}

	inst := resp.Items[0]
	r.cacheInstrument(inst)
	return inst.InstrumentID, inst.Symbol, nil
}

func (r *Resolver) ResolveMultiple(symbols []string) ([]int, []string, error) {
	ids := make([]int, 0, len(symbols))
	names := make([]string, 0, len(symbols))
	for _, s := range symbols {
		id, name, err := r.Resolve(s)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, id)
		names = append(names, name)
	}
	return ids, names, nil
}

func (r *Resolver) cacheInstrument(inst api.Instrument) {
	key := strings.ToUpper(inst.Symbol)
	r.cache[key] = cachedEntry{
		InstrumentID: inst.InstrumentID,
		Symbol:       inst.Symbol,
		DisplayName:  inst.DisplayName,
		CachedAt:     time.Now().Unix(),
	}
	r.dirty = true
	r.saveCache()
}

func (r *Resolver) CachedSymbols() []string {
	symbols := make([]string, 0, len(r.cache))
	for _, entry := range r.cache {
		symbols = append(symbols, entry.Symbol)
	}
	return symbols
}

func cachePath() string {
	return filepath.Join(config.CacheDir(), "instruments.json")
}

func (r *Resolver) loadCache() {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &r.cache)
}

func (r *Resolver) saveCache() {
	if !r.dirty {
		return
	}
	dir := config.CacheDir()
	_ = os.MkdirAll(dir, 0o700)

	data, err := json.MarshalIndent(r.cache, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(cachePath(), data, 0o600)
	r.dirty = false
}
