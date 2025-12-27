-- Test USE database switching
-- Run with: hive-duck -f test/scripts/20_use_database_switch.sql --config test/databases.yaml

-- Start in analytics (default from config)
SELECT 'Current: analytics' AS info;

-- Create/insert in analytics
CREATE TABLE IF NOT EXISTS log_entries (id INTEGER, msg VARCHAR);
INSERT INTO log_entries VALUES (1, 'Started');

-- Switch to warehouse
USE warehouse;
SELECT 'Switched to: warehouse' AS info;

-- Create/insert in warehouse
CREATE TABLE IF NOT EXISTS inventory (item_id INTEGER, quantity INTEGER);
INSERT INTO inventory VALUES (100, 50);

-- Switch to staging (in-memory)
USE staging;
SELECT 'Switched to: staging (in-memory)' AS info;

-- Create/insert in staging
CREATE TABLE temp_data (val INTEGER);
INSERT INTO temp_data VALUES (999);

-- Verify all databases have their data
SELECT 'analytics.log_entries' AS source, * FROM analytics.log_entries;
SELECT 'warehouse.inventory' AS source, * FROM warehouse.inventory;
SELECT 'staging.temp_data' AS source, * FROM staging.temp_data;


