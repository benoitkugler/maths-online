-- add a boolean field to FunctionsGraphBlock
BEGIN;
--
--

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionsGraphBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(KEY IN ('FunctionExprs', 'FunctionVariations', 'SequenceExprs', 'Areas', 'Points', 'ShowGrid', 'ShowOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'FunctionExprs')
        AND gomacro_validate_json_array_ques_VariationTableBlock (data -> 'FunctionVariations')
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'SequenceExprs')
        AND gomacro_validate_json_array_ques_FunctionArea (data -> 'Areas')
        AND gomacro_validate_json_array_ques_FunctionPoint (data -> 'Points')
        AND gomacro_validate_json_boolean (data -> 'ShowGrid')
        AND gomacro_validate_json_boolean (data -> 'ShowOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
-- 
-- 
CREATE OR REPLACE FUNCTION __migration_functions (tree jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN tree || '{"ShowGrid": true, "ShowOrigin": true}'::jsonb;
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
                CASE WHEN value ->> 'Kind' = 'FunctionsGraphBlock' THEN
                    jsonb_set(value, '{Data}', __migration_functions (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(enonce)), '[]'),
    correction = coalesce((
        SELECT
            jsonb_agg(
                CASE WHEN value ->> 'Kind' = 'FunctionsGraphBlock' THEN
                    jsonb_set(value, '{Data}', __migration_functions (value -> 'Data'))
            ELSE
                value
                END)
            FROM jsonb_array_elements(correction)), '[]');
--
--

DROP FUNCTION __migration_functions (jsonb);
COMMIT;

