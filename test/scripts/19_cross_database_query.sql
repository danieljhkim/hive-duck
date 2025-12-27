-- Test cross-database queries
-- Run with: hive-duck -f test/scripts/19_cross_database_query.sql --config test/databases.yaml
-- Requires: 18_database_mapping.sql to have been run first

-- Cross-database join: analytics.events JOIN warehouse.users
SELECT 
  e.event_id,
  e.event_type,
  u.name AS user_name
FROM analytics.events e
JOIN warehouse.users u ON e.user_id = u.user_id
ORDER BY e.event_id;

-- Aggregation across databases
SELECT 
  u.name,
  COUNT(e.event_id) AS event_count,
  COUNT(CASE WHEN e.event_type = 'click' THEN 1 END) AS clicks
FROM warehouse.users u
LEFT JOIN analytics.events e ON u.user_id = e.user_id
GROUP BY u.user_id, u.name
ORDER BY event_count DESC;


