# hive-duck

A local-development replacement for `hive -e` and `hive -f` backed by DuckDB.

## Install

```bash
go install github.com/danieljhkim/hive-duck@latest
```

Or build from source:
```bash
make build
```

## Quick Start

```bash
# Execute SQL string
hive-duck -e "SELECT 1"

# Execute SQL file
hive-duck -f script.sql

# With variables
hive-duck -e "SELECT '${hivevar:ds}'" --hivevar ds=2025-01-15

# Database mapping (YAML config)
hive-duck -f query.sql --config databases.yaml

# Output formats
hive-duck -e "SELECT * FROM table" --output json
```

## Features

- **Hive-compatible CLI**: `-e` and `-f` flags matching Hive behavior
- **Variable substitution**: `${hivevar:...}`, `${hiveconf:...}`, `${env:...}`
- **Hive statement rewrites**: `SET` and `USE` statements handled automatically
- **Database mapping**: YAML config maps Hive databases to DuckDB files
- **Output formats**: `table` (default), `csv`, `tsv`, `json`
- **Unsupported detection**: Warn or fail on unsupported Hive statements
- **DuckDB extensions**: Load extensions via `--ext avro,httpfs,json`

## Database Mapping

Create `databases.yaml`:
```yaml
databases:
  analytics: ./data/analytics.duckdb
  warehouse: ./data/warehouse.duckdb
default: analytics
```

Use with `--config databases.yaml` to enable cross-database queries.

## Flags

| Flag | Description |
|------|-------------|
| `-e, --execute` | SQL string to execute |
| `-f, --file` | SQL file to execute |
| `--config` | Path to databases.yaml |
| `--database` | DuckDB file path or `:memory:` |
| `--output` | Output format: `table`, `csv`, `tsv`, `json` |
| `--dry-run` | Print rewritten SQL without executing |
| `--fail-on-unsupported` | Fail if unsupported Hive statements detected |
| `--hivevar`, `--hiveconf` | Pass variables (repeatable) |
| `--ext` | Comma-separated DuckDB extensions |

## Development

```bash
make test-all    # Run all tests
make test-golden # Run golden output tests
make ci          # Full CI check (format, lint, test)
```

## License

MIT
