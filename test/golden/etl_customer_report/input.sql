-- ETL Customer Report Test
-- Multi-table JOIN with window functions and variable substitution

SET mapred.job.name=customer_report;

-- Create customers table
CREATE TABLE customers (
    customer_id INTEGER,
    name VARCHAR,
    segment VARCHAR,
    created_date DATE
);

-- Create orders table
CREATE TABLE orders (
    order_id INTEGER,
    customer_id INTEGER,
    order_date DATE,
    amount DOUBLE
);

-- Insert customer data
INSERT INTO customers VALUES
    (1, 'Alice Corp', 'Enterprise', '2024-01-15'),
    (2, 'Bob LLC', 'SMB', '2024-03-20'),
    (3, 'Charlie Inc', 'Enterprise', '2024-02-10'),
    (4, 'Delta Co', 'SMB', '2024-06-01');

-- Insert order data
INSERT INTO orders VALUES
    (101, 1, '2025-01-10', 5000),
    (102, 1, '2025-01-15', 3000),
    (103, 2, '2025-01-12', 1500),
    (104, 3, '2025-01-08', 8000),
    (105, 3, '2025-01-20', 2000),
    (106, 1, '2025-01-22', 4500);

-- Customer report with window functions
SELECT 
    c.name,
    c.segment,
    o.order_date,
    o.amount,
    SUM(o.amount) OVER (PARTITION BY c.customer_id ORDER BY o.order_date) AS running_total,
    ROW_NUMBER() OVER (PARTITION BY c.customer_id ORDER BY o.order_date) AS order_seq
FROM customers c
JOIN orders o ON c.customer_id = o.customer_id
ORDER BY c.name, o.order_date;

