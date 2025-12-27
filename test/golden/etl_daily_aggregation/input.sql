-- ETL Daily Aggregation Test
-- Simulates a daily ETL job that aggregates event data

SET hive.exec.dynamic.partition=true;
SET hive.exec.dynamic.partition.mode=nonstrict;

-- Create source table
CREATE TABLE raw_events (
    event_id INTEGER,
    user_id INTEGER,
    event_type VARCHAR,
    event_timestamp TIMESTAMP,
    amount DOUBLE
);

-- Insert sample data
INSERT INTO raw_events VALUES
    (1, 101, 'purchase', '2025-01-15 10:30:00', 99.99),
    (2, 102, 'view', '2025-01-15 11:00:00', 0),
    (3, 101, 'purchase', '2025-01-15 14:20:00', 49.50),
    (4, 103, 'purchase', '2025-01-15 16:45:00', 199.00),
    (5, 102, 'purchase', '2025-01-15 17:00:00', 75.25),
    (6, 101, 'view', '2025-01-15 18:00:00', 0);

-- Create aggregated table
CREATE TABLE daily_user_summary (
    user_id INTEGER,
    event_date DATE,
    total_purchases INTEGER,
    total_views INTEGER,
    total_amount DOUBLE
);

-- Aggregate and insert
INSERT INTO daily_user_summary
SELECT 
    user_id,
    CAST(event_timestamp AS DATE) AS event_date,
    COUNT(CASE WHEN event_type = 'purchase' THEN 1 END) AS total_purchases,
    COUNT(CASE WHEN event_type = 'view' THEN 1 END) AS total_views,
    SUM(CASE WHEN event_type = 'purchase' THEN amount ELSE 0 END) AS total_amount
FROM raw_events
GROUP BY user_id, CAST(event_timestamp AS DATE);

-- Output the results
SELECT * FROM daily_user_summary ORDER BY total_amount DESC;

