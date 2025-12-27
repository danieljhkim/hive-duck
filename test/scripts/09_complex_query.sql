-- Complex query with aggregations and joins
CREATE TABLE orders (
  order_id INTEGER,
  customer_id INTEGER,
  product_id INTEGER,
  quantity INTEGER,
  price DOUBLE
);

CREATE TABLE customers (
  customer_id INTEGER,
  name VARCHAR,
  city VARCHAR
);

INSERT INTO orders VALUES
  (1, 101, 1, 2, 10.0),
  (2, 101, 2, 1, 15.0),
  (3, 102, 1, 3, 10.0),
  (4, 103, 2, 2, 15.0);

INSERT INTO customers VALUES
  (101, 'Alice', 'New York'),
  (102, 'Bob', 'San Francisco'),
  (103, 'Charlie', 'New York');

SELECT 
  c.name,
  c.city,
  COUNT(o.order_id) AS order_count,
  SUM(o.quantity * o.price) AS total_spent
FROM customers c
LEFT JOIN orders o ON c.customer_id = o.customer_id
GROUP BY c.customer_id, c.name, c.city
ORDER BY total_spent DESC;

