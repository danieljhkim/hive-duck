-- Test variable substitution with hivevar
SELECT '${hivevar:ds}' AS date_string, '${hivevar:table_name}' AS table_name;

-- Test hiveconf variables
SELECT '${hiveconf:output_path}' AS output_path;

-- Test environment variables
SELECT '${env:USER}' AS username, '${env:HOME}' AS home_dir;

