-- Test script for --dry-run mode
-- Run with: hive-duck -f test/scripts/16_dry_run_test.sql --dry-run
-- Should print rewritten SQL without executing

SET hive.exec.dynamic.partition=true;

USE my_database;

SELECT 1 AS test;

CREATE TABLE test_table (id INTEGER);

INSERT INTO test_table VALUES (1);

SELECT * FROM test_table;


