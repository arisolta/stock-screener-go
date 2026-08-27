package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stock-screener/internal/yfinance"
)

const (
	defaultCacheDir = ".cache/capital_screener"
	defaultTTL      = 24 * time.Hour
)

type Cache struct {
	dir string
	ttl time.Duration
}

func New(dir string, ttl time.Duration) *Cache {
	if dir == "" {
		dir = defaultCacheDir
	}
	if ttl == 0 {
		ttl = defaultTTL
	}
	return &Cache{
		dir: dir,
		ttl: ttl,
	}
}

func (c *Cache) filePath(ticker string) string {
	safe := strings.ReplaceAll(ticker, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	return filepath.Join(c.dir, safe+".json")
}

func (c *Cache) Get(ticker string) (*yfinance.RawFinancialData, bool) {
	fp := c.filePath(ticker)
	info, err := os.Stat(fp)
	if err != nil {
		return nil, false
	}

	if time.Since(info.ModTime()) > c.ttl {
		return nil, false
	}

	bytes, err := os.ReadFile(fp)
	if err != nil {
		return nil, false
	}

	var data yfinance.RawFinancialData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, false
	}

	return &data, true
}

func (c *Cache) Set(ticker string, data *yfinance.RawFinancialData) error {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}

	fp := c.filePath(ticker)
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fp, bytes, 0644)
}
