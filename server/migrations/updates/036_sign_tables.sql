-- v1.5
-- add support for multi line signe tables

BEGIN;
--
--

CREATE OR REPLACE FUNCTION __migration_sign_table (sign_table jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_set(sign_table, '{Functions}', jsonb_build_array(jsonb_build_object('Label', sign_table -> 'Label', 'FxSymbols', sign_table -> 'FxSymbols', 'Signs', sign_table -> 'Signs'))) - 'Label' - 'Signs' - 'FxSymbols';
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
--
--
--

UPDATE
    questions
SET
    enonce = coalesce((
        SELECT
            jsonb_agg(
                CASE WHEN value ->> 'Kind' = 'SignTableBlock' THEN
                    jsonb_set(value, '{Data}', __migration_sign_table (value -> 'Data'))
                WHEN value ->> 'Kind' = 'SignTableFieldBlock' THEN
                    jsonb_set(value, '{Data, Answer}', __migration_sign_table (value -> 'Data' -> 'Answer'))
                ELSE
                    value
                END)
            FROM jsonb_array_elements(enonce)), '[]'),
    correction = coalesce((
        SELECT
            jsonb_agg(
                CASE WHEN value ->> 'Kind' = 'SignTableBlock' THEN
                    jsonb_set(value, '{Data}', __migration_sign_table (value -> 'Data'))
                WHEN value ->> 'Kind' = 'SignTableFieldBlock' THEN
                    jsonb_set(value, '{Data, Answer}', __migration_sign_table (value -> 'Data' -> 'Answer'))
                ELSE
                    value
                END)
            FROM jsonb_array_elements(correction)), '[]');
--
--

DROP FUNCTION __migration_sign_table (jsonb);
COMMIT;

