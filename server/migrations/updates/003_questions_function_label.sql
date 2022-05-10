-- v0.2.5
-- add Label to VariationTableBlock (18), SignTableBlock (13), FunctionVariationGraphBlock (9),
-- VariationTableFieldBlock (19)
-- defaulting to f

UPDATE
    questions
SET
    enonce = (
        SELECT
            jsonb_agg(
                CASE WHEN (value -> 'Kind')::int = 9 THEN
                    jsonb_set(value, '{Data, Label}', '"f"')
                WHEN (value -> 'Kind')::int = 13 THEN
                    jsonb_set(value, '{Data, Label}', '"f"')
                WHEN (value -> 'Kind')::int = 18 THEN
                    jsonb_set(value, '{Data, Label}', '"f"')
                WHEN (value -> 'Kind')::int = 19 THEN
                    jsonb_set(value, '{Data, Answer, Label}', '"f"')
                ELSE
                    value
                END)
        FROM
            jsonb_array_elements(enonce));

