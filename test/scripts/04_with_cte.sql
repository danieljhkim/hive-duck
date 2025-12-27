-- Test WITH clause (CTE)
WITH numbers AS (
  SELECT 1 AS n UNION ALL
  SELECT 2 UNION ALL
  SELECT 3 UNION ALL
  SELECT 4 UNION ALL
  SELECT 5
)
SELECT n, n * n AS squared, n * n * n AS cubed
FROM numbers
ORDER BY n;

