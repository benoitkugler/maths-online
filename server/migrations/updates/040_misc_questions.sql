-- add IsDiscrete field to function graph
BEGIN;
--
--

CREATE OR REPLACE FUNCTION __migration_point_func (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN data || '{ "IsDiscrete" : false}'::jsonb;
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
                CASE WHEN value ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
                    jsonb_set(value, '{Data}', __migration_point_func (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(enonce)), '[]'),
    correction = coalesce((
        SELECT
            jsonb_agg(
                CASE WHEN value ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
                    jsonb_set(value, '{Data}', __migration_point_func (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(correction)), '[]');
--
--

DROP FUNCTION __migration_point_func (jsonb);
COMMIT;

