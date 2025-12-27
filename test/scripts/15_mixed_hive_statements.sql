-- Test mixed Hive and standard SQL statements
-- Combines SET, USE, and regular SQL

SET hive.mapred.mode=strict;
SET hive.exec.parallel=true;

USE analytics;

CREATE TABLE events (
  event_id INTEGER,
  event_type VARCHAR,
  event_date DATE
);

INSERT INTO events VALUES
  (1, 'click', '2025-01-01'),
  (2, 'view', '2025-01-01'),
  (3, 'click', '2025-01-02');

SELECT 
  event_type,
  COUNT(*) AS event_count
FROM events
GROUP BY event_type
ORDER BY event_count DESC;

