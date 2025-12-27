-- Parameterized query using variables
-- Expected usage:
-- ./hive-duck -f test/scripts/10_parameterized_query.sql \
--   --hivevar start_date=2025-01-01 \
--   --hivevar end_date=2025-01-31 \
--   --hivevar min_amount=100.0

SELECT 
  '${hivevar:start_date}' AS start_date,
  '${hivevar:end_date}' AS end_date,
  ${hivevar:min_amount} AS min_amount,
  'Filtering between ${hivevar:start_date} and ${hivevar:end_date}' AS description;

