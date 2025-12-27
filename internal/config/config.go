package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HiveConf   map[string]string
	HiveVar    map[string]string
	Env        map[string]string
	StrictVars bool
}

func FromFlags(hiveconf []string, hivevar []string) (*Config, error) {
	cfg := &Config{
		HiveConf: make(map[string]string),
		HiveVar:  make(map[string]string),
		Env:      make(map[string]string),
	}
	for _, kv := range hiveconf {
		k, v, err := parseKV(kv)
		if err != nil {
			return nil, fmt.Errorf("--hiveconf: %w", err)
		}
		cfg.HiveConf[k] = v
	}
	for _, kv := range hivevar {
		k, v, err := parseKV(kv)
		if err != nil {
			return nil, fmt.Errorf("--hivevar: %w", err)
		}
		cfg.HiveVar[k] = v
	}

	// Snapshot env for deterministic substitution (optional but useful).
	for _, e := range os.Environ() {
		if i := strings.IndexByte(e, '='); i >= 0 {
			cfg.Env[e[:i]] = e[i+1:]
		}
	}
	return cfg, nil
}

func parseKV(s string) (string, string, error) {
	i := strings.IndexByte(s, '=')
	if i <= 0 {
		return "", "", fmt.Errorf("expected k=v, got %q", s)
	}
	return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:]), nil
}
