# Test Scripts

This directory contains test SQL scripts for hive-duck.

## Script Descriptions

- `01_basic_select.sql` - Simple SELECT query
- `02_multiple_statements.sql` - Multiple statements in one file
- `03_variable_substitution.sql` - Tests `${hivevar:...}`, `${hiveconf:...}`, `${env:...}`
- `04_with_cte.sql` - Common Table Expressions (WITH clauses)
- `05_create_table.sql` - CREATE TABLE and INSERT statements
- `06_comments_and_quotes.sql` - Comment handling and quote escaping
- `07_show_describe.sql` - SHOW and DESCRIBE commands
- `08_pragma.sql` - PRAGMA commands (DuckDB-specific)
- `09_complex_query.sql` - Complex query with joins and aggregations
- `10_parameterized_query.sql` - Parameterized query using variables
- `11_empty_and_whitespace.sql` - Empty statements and whitespace handling
- `12_nested_quotes.sql` - Nested quotes and complex string handling
- `13_set_statements.sql` - Hive SET statement handling
- `14_use_database.sql` - USE database/schema handling (legacy mode)
- `15_mixed_hive_statements.sql` - Mixed Hive and standard SQL
- `16_dry_run_test.sql` - Test script for --dry-run mode
- `17_output_formats.sql` - Test different output formats
- `18_database_mapping.sql` - Database mapping with config file
- `19_cross_database_query.sql` - Cross-database JOIN queries
- `20_use_database_switch.sql` - USE database switching with persistence

## Running Tests

### Basic test
```bash
hive-duck -f test/scripts/01_basic_select.sql
```

### With variables
```bash
hive-duck -f test/scripts/03_variable_substitution.sql \
  --hivevar ds=2025-01-15 \
  --hivevar table_name=my_table \
  --hiveconf output_path=/tmp/output
```

### Parameterized query
```bash
hive-duck -f test/scripts/10_parameterized_query.sql \
  --hivevar start_date=2025-01-01 \
  --hivevar end_date=2025-01-31 \
  --hivevar min_amount=100.0
```

### With database file
```bash
hive-duck -f test/scripts/05_create_table.sql \
  --database test/data/test.duckdb
```

### Test SET statements
```bash
hive-duck -f test/scripts/13_set_statements.sql
```

### Test USE database
```bash
hive-duck -f test/scripts/14_use_database.sql
```

### Dry-run mode (preview rewritten SQL)
```bash
hive-duck -f test/scripts/16_dry_run_test.sql --dry-run
```

### Output formats
```bash
# Table format (default)
hive-duck -f test/scripts/17_output_formats.sql --output table

# CSV format
hive-duck -f test/scripts/17_output_formats.sql --output csv

# TSV format
hive-duck -f test/scripts/17_output_formats.sql --output tsv

# JSON format
hive-duck -f test/scripts/17_output_formats.sql --output json
```

## Database Mapping Tests

These tests use a YAML config file to map Hive database names to DuckDB files.

### Setup and basic mapping
```bash
# Creates analytics and warehouse databases with test data
hive-duck -f test/scripts/18_database_mapping.sql --config test/databases.yaml
```

### Cross-database queries
```bash
# JOINs data across analytics and warehouse databases
hive-duck -f test/scripts/19_cross_database_query.sql --config test/databases.yaml
```

### USE database switching
```bash
# Tests switching between databases with USE statements
hive-duck -f test/scripts/20_use_database_switch.sql --config test/databases.yaml
```

### Config file format (test/databases.yaml)
```yaml
databases:
  analytics: ./data/analytics.duckdb
  warehouse: ./data/warehouse.duckdb
  staging: ":memory:"
default: analytics
```
