# hive-duck

**hive-duck** is a Hive-compatible command-line tool that lets you run `hive -e` and `hive -f` style SQL locally using DuckDB instead of Hadoop and Hive.

It is designed for **fast local development and testing of Hive SQL** without standing up HDFS, YARN, or the Hive Metastore. Existing scripts can usually be run unchanged, while DuckDB provides a single-node, high-performance analytical execution engine.

Typical use cases:
- Develop and debug Hive SQL locally
- Run CI checks for SQL pipelines without Hadoop
- Replace local Hive-on-MapReduce or Spark SQL setups
- Iterate quickly on Parquet-, Avro-, or CSV-based datasets

hive-duck focuses on **CLI compatibility and correctness**, not distributed execution or performance parity with Spark.

## Features

- **Hive-compatible execution**  
  Supports `-e` (inline SQL) and `-f` (SQL files) with behavior aligned to the Hive CLI.

- **Variable substitution**  
  Compatible with Hive-style variables: `${hivevar:...}`, `${hiveconf:...}`, and `${env:...}`.

- **Hive statement handling**  
  Common Hive statements such as `SET` and `USE` are rewritten or handled automatically for local execution.

- **DuckDB-backed execution**  
  Runs SQL using DuckDB’s in-process analytical engine for fast, single-node execution.

- **Database mapping via YAML**  
  Map Hive databases to DuckDB database files using a simple configuration file.

- **Multiple output formats**  
  Render query results as `table` (default), `csv`, `tsv`, or `json`.

- **Extension support**  
  Load DuckDB extensions (e.g. `avro`, `httpfs`, `json`) via a single flag.

- **Compatibility checks**  
  Detect unsupported Hive statements and optionally fail fast for CI use cases.


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
make build        # Build the binary
make test-all     # Run unit and integration tests
make test-golden  # Run golden SQL output tests
make ci           # Format, lint, and test
```

## License

MIT
