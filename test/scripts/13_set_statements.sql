-- Test SET statement handling
-- SET statements should be captured but not executed

SET hive.exec.dynamic.partition=true;
SET hive.exec.dynamic.partition.mode=nonstrict;
SET mapred.job.name="test_job";

-- This SELECT should still work after SET statements
SELECT 'SET statements captured successfully' AS result;


