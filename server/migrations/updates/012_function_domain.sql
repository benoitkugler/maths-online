-- v0.5.0
-- uses expressions instead of fixed values for function definition

CREATE OR REPLACE FUNCTION __migration_intArrayToStringArray (elements jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN coalesce((
        SELECT
            (jsonb_agg(to_jsonb ((value::int)::text)))
    FROM jsonb_array_elements(elements)), '[]');
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_functionRangeArray (elements jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN coalesce((
        SELECT
            (jsonb_agg(jsonb_set(jsonb_set(value - 'Range', '{From}', to_jsonb (value -> 'Range' ->> 0)), '{To}', to_jsonb (value -> 'Range' ->> 1))))
        FROM jsonb_array_elements(elements)), '[]');
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', (
            SELECT
                jsonb_agg(
                    CASE WHEN (value -> 'Kind')::int = 8 THEN
                        jsonb_set(value, '{Data, XGrid}', __migration_intArrayToStringArray (value -> 'Data' -> 'XGrid'))
                    WHEN (value -> 'Kind')::int = 7 THEN
                        jsonb_set(value, '{Data, Functions}', __migration_functionRangeArray (value -> 'Data' -> 'Functions'))
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')));

-- Put the CONSTRAINT back, after updating the definitions
ALTER TABLE questions
    ADD CONSTRAINT question_structgen_validate_json_exe_questionpage CHECK (structgen_validate_json_exe_QuestionPage (page));

