-- ETL Cross-Database Test
-- Cross-database JOIN using config mapping

SET hive.exec.mode.local.auto=true;

USE analytics;

-- Create events table in analytics
CREATE TABLE IF NOT EXISTS page_views (
    view_id INTEGER,
    user_id INTEGER,
    page VARCHAR,
    view_time TIMESTAMP
);

INSERT INTO page_views VALUES
    (1, 201, '/home', '2025-01-15 09:00:00'),
    (2, 202, '/products', '2025-01-15 09:05:00'),
    (3, 201, '/products', '2025-01-15 09:10:00'),
    (4, 203, '/home', '2025-01-15 09:15:00'),
    (5, 201, '/checkout', '2025-01-15 09:20:00');

USE warehouse;

-- Create users table in warehouse
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id INTEGER,
    username VARCHAR,
    tier VARCHAR
);

INSERT INTO user_profiles VALUES
    (201, 'alice', 'premium'),
    (202, 'bob', 'basic'),
    (203, 'charlie', 'premium');

-- Cross-database query: join analytics.page_views with warehouse.user_profiles
SELECT 
    u.username,
    u.tier,
    COUNT(p.view_id) AS view_count,
    COUNT(DISTINCT p.page) AS unique_pages
FROM warehouse.user_profiles u
LEFT JOIN analytics.page_views p ON u.user_id = p.user_id
GROUP BY u.user_id, u.username, u.tier
ORDER BY view_count DESC, u.username;

