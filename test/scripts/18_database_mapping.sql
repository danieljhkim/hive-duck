-- Test database mapping feature
-- Run with: hive-duck -f test/scripts/18_database_mapping.sql --config test/databases.yaml

-- Create table in analytics database (default)
CREATE TABLE IF NOT EXISTS events (
  event_id INTEGER,
  event_type VARCHAR,
  user_id INTEGER
);

INSERT INTO events VALUES (1, 'click', 101), (2, 'view', 102), (3, 'click', 101);

-- Switch to warehouse database
USE warehouse;

CREATE TABLE IF NOT EXISTS users (
  user_id INTEGER,
  name VARCHAR
);

INSERT INTO users VALUES (101, 'Alice'), (102, 'Bob');

-- Query from warehouse
SELECT * FROM users;


