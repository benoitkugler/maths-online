-- fix simple indices in QCM fields
-- first check that there is no complex expression yet

SELECT DISTINCT
    (
        SELECT
            jsonb_agg(value -> 'Data' -> 'Answer')
        FROM
            jsonb_array_elements(enonce)
        WHERE (value -> 'Kind')::int = 12)
FROM
    questions;

UPDATE
    questions
SET
    enonce = (
        SELECT
            jsonb_agg(
                CASE WHEN (value -> 'Kind')::int = 12 THEN
                    CASE WHEN value -> 'Data' ->> 'Answer' = '0' THEN
                        jsonb_set(value, '{Data, Answer}', '"1"')
                    WHEN value -> 'Data' ->> 'Answer' = '1' THEN
                        jsonb_set(value, '{Data, Answer}', '"2"')
                    ELSE
                        value
                    END
                ELSE
                    value
                END)
        FROM
            jsonb_array_elements(enonce));

-- check if the fix is OK
SELECT DISTINCT
    (
        SELECT
            jsonb_agg(value -> 'Data' -> 'Answer')
        FROM
            jsonb_array_elements(enonce)
        WHERE (value -> 'Kind')::int = 12)
FROM
    questions;

