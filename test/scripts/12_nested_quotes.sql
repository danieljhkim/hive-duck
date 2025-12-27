-- Test nested quotes and complex string handling
SELECT 'Outer ''inner'' quote' AS nested_single;
SELECT "Outer ""inner"" quote" AS nested_double;
SELECT 'Mixed "quotes" inside' AS mixed_quotes;
SELECT 'String with ; semicolon; and ''quotes''' AS complex;

