package output

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

// Format represents the output format for query results.
type Format string

const (
	FormatTable Format = "table" // Default: aligned columns using tabwriter
	FormatCSV   Format = "csv"   // RFC 4180 CSV
	FormatTSV   Format = "tsv"   // Tab-separated values
	FormatJSON  Format = "json"  // JSON array of objects
)

// ValidFormats returns a list of valid output format names.
func ValidFormats() []string {
	return []string{string(FormatTable), string(FormatCSV), string(FormatTSV), string(FormatJSON)}
}

// ParseFormat parses a format string and returns the corresponding Format.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "table", "":
		return FormatTable, nil
	case "csv":
		return FormatCSV, nil
	case "tsv":
		return FormatTSV, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("invalid output format %q (valid: %s)", s, strings.Join(ValidFormats(), ", "))
	}
}

// PrintRows outputs query results in the specified format.
func PrintRows(rows *sql.Rows, format Format) error {
	switch format {
	case FormatTable:
		return printTable(rows)
	case FormatCSV:
		return printCSV(rows, ',')
	case FormatTSV:
		return printCSV(rows, '\t')
	case FormatJSON:
		return printJSON(rows)
	default:
		return printTable(rows)
	}
}

// printTable outputs results as aligned columns (original behavior).
func printTable(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		if err := tw.Flush(); err != nil {
			log.Printf("Failed to flush tabwriter: %v", err)
		}
	}()

	// Header
	for i, c := range cols {
		if i > 0 {
			_, _ = io.WriteString(tw, "\t")
		}
		_, _ = io.WriteString(tw, c)
	}
	_, _ = io.WriteString(tw, "\n")

	// Scan buffers
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		for i, v := range vals {
			if i > 0 {
				_, _ = io.WriteString(tw, "\t")
			}
			_, _ = io.WriteString(tw, formatValue(v))
		}
		_, _ = io.WriteString(tw, "\n")
	}
	return rows.Err()
}

// printCSV outputs results as CSV or TSV.
func printCSV(rows *sql.Rows, delimiter rune) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = delimiter
	defer w.Flush()

	// Header
	if err := w.Write(cols); err != nil {
		return err
	}

	// Scan buffers
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	record := make([]string, len(cols))
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		for i, v := range vals {
			record[i] = formatValue(v)
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return rows.Err()
}

// printJSON outputs results as a JSON array of objects.
func printJSON(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// Scan buffers
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	var results []map[string]any

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return err
		}
		row := make(map[string]any, len(cols))
		for i, col := range cols {
			row[col] = convertJSONValue(vals[i])
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// Handle empty results
	if results == nil {
		results = []map[string]any{}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

// formatValue converts a database value to a string for display.
func formatValue(v any) string {
	if v == nil {
		return "NULL"
	}
	return fmt.Sprint(v)
}

// convertJSONValue converts a database value for JSON serialization.
func convertJSONValue(v any) any {
	if v == nil {
		return nil
	}
	// Handle []byte as string for JSON
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return v
}
