-- v0.6.1
-- merge FunctionExpressions and FunctionVariation in one block

CREATE OR REPLACE FUNCTION __migration_functionExprToFunctions (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object('FunctionExprs', data -> 'Functions', 'FunctionVariations', '[]'::jsonb, 'Areas', '[]'::jsonb);
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_functionVariationToFunctions (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object('FunctionExprs', '[]'::jsonb, 'FunctionVariations', jsonb_build_array(data), 'Areas', '[]'::jsonb);
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN (value ->> 'Kind') = 'FunctionGraphBlock' THEN
                        jsonb_set(jsonb_set(value, '{Kind}', to_jsonb ('FunctionsGraphBlock'::text)), '{Data}', __migration_functionExprToFunctions (value -> 'Data'))
                    WHEN (value ->> 'Kind') = 'FunctionVariationGraphBlock' THEN
                        jsonb_set(jsonb_set(value, '{Kind}', to_jsonb ('FunctionsGraphBlock'::text)), '{Data}', __migration_functionVariationToFunctions (value -> 'Data'))
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));

ALTER TABLE questions
    ADD CONSTRAINT page_structgen_validate_json_que_QuestionPage CHECK (structgen_validate_json_que_QuestionPage (page));

