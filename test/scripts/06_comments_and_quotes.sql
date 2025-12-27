-- Test handling of comments and quotes
-- This is a single-line comment

/* This is a
   multi-line comment */

SELECT 'Single quoted string' AS single_quote;
SELECT "Double quoted string" AS double_quote;
SELECT 'String with ''escaped'' quotes' AS escaped;
SELECT 'String with semicolon; inside' AS with_semicolon;

-- Statement after comments
SELECT 42 AS answer;

