package engine

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/marcboeker/go-duckdb"

	"github.com/danieljhkim/hive-duck/internal/config"
	"github.com/danieljhkim/hive-duck/internal/output"
)

type Runner struct {
	DBPath       string
	Exts         []string
	Silent       bool
	OutputFormat output.Format
	DatabaseMap  *config.DatabaseMap // Optional: Hive DB -> DuckDB path mapping
}

func (r Runner) Run(stmts []string) error {
	// go-duckdb uses empty string for in-memory database, not ":memory:"
	dsn := r.DBPath
	if dsn == ":memory:" {
		dsn = ""
	}

	db, err := sql.Open("duckdb", dsn)
	if err != nil {
		return fmt.Errorf("open duckdb: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Extensions
	for _, ext := range r.Exts {
		if err := exec(db, fmt.Sprintf("INSTALL %s", ident(ext))); err != nil {
			return fmt.Errorf("install ext %q: %w", ext, err)
		}
		if err := exec(db, fmt.Sprintf("LOAD %s", ident(ext))); err != nil {
			return fmt.Errorf("load ext %q: %w", ext, err)
		}
	}

	// ATTACH mapped databases
	if r.DatabaseMap != nil {
		if err := r.attachDatabases(db); err != nil {
			return err
		}
	}

	for _, stmt := range stmts {
		trim := strings.TrimSpace(stmt)
		if trim == "" {
			continue
		}

		// Heuristic: print results if it looks like it returns rows
		if returnsRows(trim) {
			rows, err := db.Query(trim)
			if err != nil {
				return fmt.Errorf("query failed: %w\nSQL: %s", err, trim)
			}
			if err := output.PrintRows(rows, r.OutputFormat); err != nil {
				_ = rows.Close()
				return err
			}
			_ = rows.Close()
			continue
		}

		if err := exec(db, trim); err != nil {
			return fmt.Errorf("exec failed: %w\nSQL: %s", err, trim)
		}
	}
	return nil
}

// attachDatabases ATTACHes all mapped databases so they're available for cross-db queries.
func (r Runner) attachDatabases(db *sql.DB) error {
	for name, path := range r.DatabaseMap.Databases {
		var attachSQL string
		if path == ":memory:" {
			// For in-memory databases, ATTACH without a path
			attachSQL = fmt.Sprintf("ATTACH '' AS %s", ident(name))
		} else {
			// For file-based databases
			attachSQL = fmt.Sprintf("ATTACH %s AS %s", quoteLiteral(path), ident(name))
		}

		if err := exec(db, attachSQL); err != nil {
			return fmt.Errorf("attach database %q (%s): %w", name, path, err)
		}
	}

	// USE default database if specified
	if r.DatabaseMap.Default != "" {
		if err := exec(db, fmt.Sprintf("USE %s", ident(r.DatabaseMap.Default))); err != nil {
			return fmt.Errorf("use default database %q: %w", r.DatabaseMap.Default, err)
		}
	}

	return nil
}

func exec(db *sql.DB, stmt string) error {
	_, err := db.Exec(stmt)
	return err
}

func returnsRows(stmt string) bool {
	s := strings.ToLower(strings.TrimSpace(stmt))
	return strings.HasPrefix(s, "select") ||
		strings.HasPrefix(s, "with") ||
		strings.HasPrefix(s, "show") ||
		strings.HasPrefix(s, "describe") ||
		strings.HasPrefix(s, "pragma")
}

// ident validates/quotes an identifier for INSTALL/LOAD/ATTACH statements.
func ident(s string) string {
	// DuckDB identifiers are typically simple; keep it strict.
	for _, ch := range s {
		if ch != '_' && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') {
			// Fall back to a quoted identifier
			return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
		}
	}
	return s
}

// quoteLiteral quotes a string literal for SQL.
func quoteLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
