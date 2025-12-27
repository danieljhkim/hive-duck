-- Test handling of empty statements and whitespace

SELECT 'First' AS col1;

-- Empty line above, statement below
SELECT 'Second' AS col1;

   -- Statement with leading whitespace
   SELECT 'Third' AS col1;

SELECT 'Fourth' AS col1;    -- Trailing whitespace

