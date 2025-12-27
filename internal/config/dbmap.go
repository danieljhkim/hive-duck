package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DatabaseMap holds the mapping of Hive database names to DuckDB database paths.
type DatabaseMap struct {
	Databases map[string]string `yaml:"databases"` // db_name -> path/to/file.duckdb
	Default   string            `yaml:"default"`   // Default database to USE on startup
}

// LoadDatabaseMap loads a database mapping from a YAML file.
func LoadDatabaseMap(path string) (*DatabaseMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read database config: %w", err)
	}

	var dbMap DatabaseMap
	if err := yaml.Unmarshal(data, &dbMap); err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	if dbMap.Databases == nil {
		dbMap.Databases = make(map[string]string)
	}

	// Resolve relative paths based on config file location
	configDir := filepath.Dir(path)
	for name, dbPath := range dbMap.Databases {
		if dbPath != ":memory:" && !filepath.IsAbs(dbPath) {
			dbMap.Databases[name] = filepath.Join(configDir, dbPath)
		}
	}

	return &dbMap, nil
}

// GetDatabasePath returns the DuckDB path for a given database name.
// Returns empty string if not found.
func (m *DatabaseMap) GetDatabasePath(name string) string {
	if m == nil || m.Databases == nil {
		return ""
	}
	return m.Databases[name]
}

// HasDatabase returns true if the database name is mapped.
func (m *DatabaseMap) HasDatabase(name string) bool {
	if m == nil || m.Databases == nil {
		return false
	}
	_, ok := m.Databases[name]
	return ok
}

// DatabaseNames returns all mapped database names.
func (m *DatabaseMap) DatabaseNames() []string {
	if m == nil || m.Databases == nil {
		return nil
	}
	names := make([]string, 0, len(m.Databases))
	for name := range m.Databases {
		names = append(names, name)
	}
	return names
}
