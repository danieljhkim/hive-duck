-- ETL Incremental Load Test
-- Date-partitioned logic with hivevar parameters

SET hive.exec.parallel=true;

-- Create staging table
CREATE TABLE staging_transactions (
    txn_id INTEGER,
    account_id INTEGER,
    txn_date DATE,
    txn_type VARCHAR,
    amount DOUBLE
);

-- Create target table
CREATE TABLE daily_transactions (
    txn_id INTEGER,
    account_id INTEGER,
    txn_date DATE,
    txn_type VARCHAR,
    amount DOUBLE,
    load_date DATE
);

-- Insert staging data (simulating incremental load for a specific date)
INSERT INTO staging_transactions VALUES
    (1001, 501, '${hivevar:ds}', 'credit', 1000.00),
    (1002, 502, '${hivevar:ds}', 'debit', 250.50),
    (1003, 501, '${hivevar:ds}', 'credit', 500.00),
    (1004, 503, '${hivevar:ds}', 'debit', 75.25);

-- Incremental insert into target table
INSERT INTO daily_transactions
SELECT 
    txn_id,
    account_id,
    txn_date,
    txn_type,
    amount,
    CURRENT_DATE AS load_date
FROM staging_transactions
WHERE txn_date = '${hivevar:ds}';

-- Summary by account for the loaded date
SELECT 
    account_id,
    COUNT(*) AS txn_count,
    SUM(CASE WHEN txn_type = 'credit' THEN amount ELSE 0 END) AS total_credits,
    SUM(CASE WHEN txn_type = 'debit' THEN amount ELSE 0 END) AS total_debits
FROM daily_transactions
WHERE txn_date = '${hivevar:ds}'
GROUP BY account_id
ORDER BY account_id;

