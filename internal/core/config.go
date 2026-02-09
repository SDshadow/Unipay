package core

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	data map[string]map[string]any
	mu   sync.RWMutex
}

var globalCfg *Config

func LoadFromBytes(b []byte) (*Config, error) {
	var raw map[string]map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &Config{data: raw}, nil
}

func LoadFromFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cfg, err := LoadFromBytes(b)
	if err != nil {
		return err
	}
	globalCfg = cfg
	return nil
}

func (c *Config) Get(channel, tenant string) map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	t, ok := c.data[channel][tenant]
	if !ok {
		return nil
	}
	return t.(map[string]any)
}

func Get(channel, tenant string) map[string]any {
	if globalCfg == nil {
		return nil
	}
	return globalCfg.Get(channel, tenant)
}
