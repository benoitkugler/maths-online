-- wrap TreeFieldBlock into a {Answer: XXX} object
BEGIN;
--
--

CREATE OR REPLACE FUNCTION __migration_trees (tree jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object('Answer', tree);
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
                CASE WHEN value ->> 'Kind' = 'TreeFieldBlock' THEN
                    jsonb_set(value, '{Data}', __migration_trees (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(enonce)), '[]'),
    correction = coalesce((
        SELECT
            jsonb_agg(
                CASE WHEN value ->> 'Kind' = 'TreeFieldBlock' THEN
                    jsonb_set(value, '{Data}', __migration_trees (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(correction)), '[]');
--
--

DROP FUNCTION __migration_trees (jsonb);
COMMIT;

