package preprocess

import (
	"fmt"
	"strings"
)

func SplitStatements(sql string) ([]string, error) {
	var (
		stmts []string
		buf   strings.Builder

		inSQuote bool
		inDQuote bool
		inLineC  bool
		inBlockC bool
	)

	r := []rune(sql)
	for i := 0; i < len(r); i++ {
		ch := r[i]

		// Comment handling
		if inLineC {
			if ch == '\n' {
				inLineC = false
				buf.WriteRune(ch)
			}
			continue
		}
		if inBlockC {
			if ch == '*' && i+1 < len(r) && r[i+1] == '/' {
				inBlockC = false
				i++
			}
			continue
		}

		// Start comments (only when not in quotes)
		if !inSQuote && !inDQuote {
			if ch == '-' && i+1 < len(r) && r[i+1] == '-' {
				inLineC = true
				i++
				continue
			}
			if ch == '/' && i+1 < len(r) && r[i+1] == '*' {
				inBlockC = true
				i++
				continue
			}
		}

		// Quote toggles (handle escaped '' inside single quotes)
		if !inDQuote && ch == '\'' {
			// If already in single quote and next is also ', treat as escaped quote
			if inSQuote && i+1 < len(r) && r[i+1] == '\'' {
				buf.WriteRune(ch)
				buf.WriteRune(r[i+1])
				i++
				continue
			}
			inSQuote = !inSQuote
			buf.WriteRune(ch)
			continue
		}
		if !inSQuote && ch == '"' {
			inDQuote = !inDQuote
			buf.WriteRune(ch)
			continue
		}

		// Statement split
		if ch == ';' && !inSQuote && !inDQuote {
			stmt := strings.TrimSpace(buf.String())
			buf.Reset()
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
			continue
		}

		buf.WriteRune(ch)
	}

	if inSQuote || inDQuote || inBlockC {
		return nil, fmt.Errorf("unterminated quote or comment in SQL input")
	}

	last := strings.TrimSpace(buf.String())
	if last != "" {
		stmts = append(stmts, last)
	}
	return stmts, nil
}
