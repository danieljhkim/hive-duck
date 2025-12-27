-- Test different output formats
-- Run with different --output options:
--   hive-duck -f test/scripts/17_output_formats.sql --output table
--   hive-duck -f test/scripts/17_output_formats.sql --output csv
--   hive-duck -f test/scripts/17_output_formats.sql --output tsv
--   hive-duck -f test/scripts/17_output_formats.sql --output json

SELECT 
  1 AS id,
  'Alice' AS name,
  25.5 AS score,
  DATE '2025-01-15' AS created
UNION ALL
SELECT 2, 'Bob', 30.2, DATE '2025-01-16'
UNION ALL
SELECT 3, 'Charlie', NULL, DATE '2025-01-17'
ORDER BY id;


