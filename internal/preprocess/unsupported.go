package preprocess

import (
	"regexp"
	"strings"
)

// UnsupportedResult represents a detected unsupported Hive statement.
type UnsupportedResult struct {
	Statement string // Original statement (truncated for display)
	Keyword   string // Detected unsupported keyword/pattern
	Reason    string // Why it's unsupported
}

// Unsupported statement patterns with descriptions
var unsupportedPatterns = []struct {
	pattern *regexp.Regexp
	keyword string
	reason  string
}{
	// Data Loading
	{
		regexp.MustCompile(`(?i)^\s*LOAD\s+DATA`),
		"LOAD DATA",
		"Use DuckDB's read_csv/read_parquet/COPY instead",
	},
	{
		regexp.MustCompile(`(?i)^\s*EXPORT\s+TABLE`),
		"EXPORT TABLE",
		"Use DuckDB's COPY TO instead",
	},
	{
		regexp.MustCompile(`(?i)^\s*IMPORT\s+TABLE`),
		"IMPORT TABLE",
		"Use DuckDB's COPY FROM instead",
	},

	// UDF/Transform
	{
		regexp.MustCompile(`(?i)^\s*ADD\s+JAR`),
		"ADD JAR",
		"Java UDFs not supported; use DuckDB extensions or native functions",
	},
	{
		regexp.MustCompile(`(?i)^\s*ADD\s+FILE`),
		"ADD FILE",
		"Distributed file adding not supported",
	},
	{
		regexp.MustCompile(`(?i)^\s*CREATE\s+(TEMPORARY\s+)?FUNCTION`),
		"CREATE FUNCTION",
		"Use DuckDB's CREATE MACRO or native functions",
	},
	{
		regexp.MustCompile(`(?i)\bTRANSFORM\s*\(`),
		"TRANSFORM",
		"Hive TRANSFORM with external scripts not supported",
	},
	{
		regexp.MustCompile(`(?i)\bMAP\s*\([^)]+\)\s+USING`),
		"MAP...USING",
		"Hive MAP transformation not supported",
	},
	{
		regexp.MustCompile(`(?i)\bREDUCE\s*\([^)]+\)\s+USING`),
		"REDUCE...USING",
		"Hive REDUCE transformation not supported",
	},

	// Hive DDL
	{
		regexp.MustCompile(`(?i)^\s*MSCK\s+REPAIR`),
		"MSCK REPAIR",
		"Partition repair not needed; DuckDB doesn't use Hive metastore",
	},
	{
		regexp.MustCompile(`(?i)^\s*ANALYZE\s+TABLE`),
		"ANALYZE TABLE",
		"Use DuckDB's ANALYZE instead",
	},
	{
		regexp.MustCompile(`(?i)^\s*TRUNCATE\s+TABLE`),
		"TRUNCATE TABLE",
		"Use DELETE FROM table or DROP TABLE + CREATE TABLE",
	},
	{
		regexp.MustCompile(`(?i)^\s*ALTER\s+TABLE\s+\S+\s+(ADD|DROP|RENAME)\s+PARTITION`),
		"ALTER TABLE...PARTITION",
		"Hive partition management not supported; use DuckDB partitioning",
	},
	{
		regexp.MustCompile(`(?i)^\s*ALTER\s+TABLE\s+\S+\s+RECOVER\s+PARTITIONS`),
		"ALTER TABLE...RECOVER PARTITIONS",
		"Partition recovery not supported",
	},

	// Hive-specific Clauses (detected within statements)
	{
		regexp.MustCompile(`(?i)\bLATERAL\s+VIEW`),
		"LATERAL VIEW",
		"Use DuckDB's UNNEST or list functions instead",
	},
	{
		regexp.MustCompile(`(?i)\bCLUSTER\s+BY\b`),
		"CLUSTER BY",
		"Use ORDER BY for sorting; clustering not supported",
	},
	{
		regexp.MustCompile(`(?i)\bDISTRIBUTE\s+BY\b`),
		"DISTRIBUTE BY",
		"Distribution hint not needed in DuckDB",
	},
	{
		regexp.MustCompile(`(?i)\bSORT\s+BY\b`),
		"SORT BY",
		"Use ORDER BY instead for deterministic sorting",
	},
	{
		regexp.MustCompile(`(?i)\bTABLESAMPLE\s*\(`),
		"TABLESAMPLE",
		"Use DuckDB's USING SAMPLE clause instead",
	},

	// Metastore operations
	{
		regexp.MustCompile(`(?i)^\s*SHOW\s+PARTITIONS`),
		"SHOW PARTITIONS",
		"Hive partitions not applicable; data is file-based",
	},
	{
		regexp.MustCompile(`(?i)^\s*SHOW\s+TBLPROPERTIES`),
		"SHOW TBLPROPERTIES",
		"Table properties not stored in Hive metastore format",
	},
	{
		regexp.MustCompile(`(?i)^\s*DESCRIBE\s+EXTENDED`),
		"DESCRIBE EXTENDED",
		"Use DESCRIBE or PRAGMA table_info instead",
	},
	{
		regexp.MustCompile(`(?i)^\s*DESCRIBE\s+FORMATTED`),
		"DESCRIBE FORMATTED",
		"Use DESCRIBE or PRAGMA table_info instead",
	},

	// Storage format hints
	{
		regexp.MustCompile(`(?i)\bSTORED\s+AS\b`),
		"STORED AS",
		"Storage format hint ignored; use read_parquet/read_csv explicitly",
	},
	{
		regexp.MustCompile(`(?i)\bROW\s+FORMAT\b`),
		"ROW FORMAT",
		"Row format specification not supported",
	},
	{
		regexp.MustCompile(`(?i)\bSERDE\b`),
		"SERDE",
		"SerDe not supported; use DuckDB's native readers",
	},
	{
		regexp.MustCompile(`(?i)\bLOCATION\s+'`),
		"LOCATION",
		"External table location; use CREATE TABLE AS or views with read_* functions",
	},

	// Locks and transactions (Hive-specific)
	{
		regexp.MustCompile(`(?i)^\s*LOCK\s+TABLE`),
		"LOCK TABLE",
		"Explicit locking not supported",
	},
	{
		regexp.MustCompile(`(?i)^\s*UNLOCK\s+TABLE`),
		"UNLOCK TABLE",
		"Explicit unlocking not supported",
	},
}

// DetectUnsupported scans statements for unsupported Hive-specific constructs.
// Returns a list of all detected issues.
func DetectUnsupported(stmts []string) []UnsupportedResult {
	var results []UnsupportedResult

	for _, stmt := range stmts {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}

		for _, p := range unsupportedPatterns {
			if p.pattern.MatchString(trimmed) {
				results = append(results, UnsupportedResult{
					Statement: truncateStatement(trimmed, 80),
					Keyword:   p.keyword,
					Reason:    p.reason,
				})
				// Don't break - a statement might match multiple patterns
			}
		}
	}

	return results
}

// HasUnsupported returns true if any unsupported statements are detected.
func HasUnsupported(stmts []string) bool {
	for _, stmt := range stmts {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		for _, p := range unsupportedPatterns {
			if p.pattern.MatchString(trimmed) {
				return true
			}
		}
	}
	return false
}

// truncateStatement shortens a statement for display.
func truncateStatement(s string, maxLen int) string {
	// Normalize whitespace
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
