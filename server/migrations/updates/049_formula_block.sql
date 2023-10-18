SELECT
    id,
    (
        SELECT
            jsonb_agg(value -> 'Data')
        FROM
            jsonb_array_elements(enonce)
        WHERE
            value ->> 'Kind' = 'FormulaBlock')
FROM
    questions;

