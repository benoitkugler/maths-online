-- v0.2.5
-- add Label to VariationTableBlock (18), SignTableBlock (13), FunctionVariationGraphBlock (9),
-- VariationTableFieldBlock (19)
-- defaulting to f
--
-- First remove the constraint (should be added back after the migration)

ALTER TABLE questions
    DROP CONSTRAINT enonce_structgen_validate_json_array_exe_Block;

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

