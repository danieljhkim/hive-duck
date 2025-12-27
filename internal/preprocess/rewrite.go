package preprocess

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/danieljhkim/hive-duck/internal/config"
)

// RewriteResult contains the output of the Hive-to-DuckDB rewrite process.
type RewriteResult struct {
	Statements    []string          // Rewritten SQL statements to execute
	CurrentSchema string            // Last schema set via USE statement
	SetVars       map[string]string // Captured SET k=v pairs (for reference)
}

// RewriteOptions configures the rewrite behavior.
type RewriteOptions struct {
	DatabaseMap *config.DatabaseMap // If set, USE statements target attached databases
}

// Regex patterns for Hive statements
var (
	// SET key=value or SET key = value (with optional quotes around value)
	setPattern = regexp.MustCompile(`(?i)^\s*SET\s+([A-Za-z0-9_.\-:]+)\s*=\s*(.*)$`)

	// USE database/schema
	usePattern = regexp.MustCompile(`(?i)^\s*USE\s+([A-Za-z0-9_]+)\s*$`)
)

// Rewrite transforms Hive SQL statements into DuckDB-compatible statements.
// - SET k=v statements are captured but not executed
// - USE db statements are rewritten based on options:
//   - With DatabaseMap: USE db (databases are pre-ATTACHed)
//   - Without DatabaseMap: CREATE SCHEMA IF NOT EXISTS + SET search_path (legacy)
func Rewrite(stmts []string, opts *RewriteOptions) (*RewriteResult, error) {
	result := &RewriteResult{
		Statements: make([]string, 0, len(stmts)),
		SetVars:    make(map[string]string),
	}

	if opts == nil {
		opts = &RewriteOptions{}
	}

	for _, stmt := range stmts {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}

		// Check for SET statement
		if matches := setPattern.FindStringSubmatch(trimmed); matches != nil {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			// Remove trailing semicolon if present
			value = strings.TrimSuffix(value, ";")
			value = strings.TrimSpace(value)
			// Remove surrounding quotes if present
			if len(value) >= 2 {
				if (value[0] == '\'' && value[len(value)-1] == '\'') ||
					(value[0] == '"' && value[len(value)-1] == '"') {
					value = value[1 : len(value)-1]
				}
			}
			result.SetVars[key] = value
			// SET statements are captured but not passed to DuckDB
			continue
		}

		// Check for USE statement
		if matches := usePattern.FindStringSubmatch(trimmed); matches != nil {
			dbName := strings.TrimSpace(matches[1])
			result.CurrentSchema = dbName

			if opts.DatabaseMap != nil {
				// Database mapping mode: verify DB is mapped, then USE it
				if !opts.DatabaseMap.HasDatabase(dbName) {
					return nil, fmt.Errorf("database %q not found in config (available: %v)",
						dbName, opts.DatabaseMap.DatabaseNames())
				}
				// Pass through as DuckDB USE (databases are pre-ATTACHed)
				result.Statements = append(result.Statements, "USE "+dbName)
			} else {
				// Legacy mode: create schema and set search_path
				result.Statements = append(result.Statements,
					"CREATE SCHEMA IF NOT EXISTS "+dbName)
				result.Statements = append(result.Statements,
					"SET search_path = '"+dbName+"'")
			}
			continue
		}

		// Pass through unchanged
		result.Statements = append(result.Statements, trimmed)
	}

	return result, nil
}

// IsHiveStatement returns true if the statement is a Hive-specific statement
// that needs special handling.
func IsHiveStatement(stmt string) bool {
	trimmed := strings.TrimSpace(stmt)
	return setPattern.MatchString(trimmed) || usePattern.MatchString(trimmed)
}
