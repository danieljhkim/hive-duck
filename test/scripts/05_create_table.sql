-- Test CREATE TABLE and INSERT
CREATE TABLE test_table (
  id INTEGER,
  name VARCHAR,
  value DOUBLE
);

INSERT INTO test_table VALUES (1, 'Alice', 10.5);
INSERT INTO test_table VALUES (2, 'Bob', 20.3);
INSERT INTO test_table VALUES (3, 'Charlie', 30.7);

SELECT * FROM test_table ORDER BY id;

