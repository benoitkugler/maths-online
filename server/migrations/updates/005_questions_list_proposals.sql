-- v0.2.5
-- generalize the proposal field to interpolated
-- automatically add $ $ delimiters to current expression

CREATE OR REPLACE FUNCTION __migration_textPartToInterpolated (textPart jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    CASE WHEN (textPart -> 'Kind')::int = 0 THEN
        -- simple text
        RETURN to_jsonb (textPart ->> 'Content');
    WHEN (textPart -> 'Kind')::int = 1 THEN
        RETURN to_jsonb (concat('$', textPart ->> 'Content', '$')); -- latex
    WHEN (textPart -> 'Kind')::int = 2 THEN
        RETURN to_jsonb (concat('&', textPart ->> 'Content', '&')); -- expression
    ELSE
        RAISE WARNING 'invalid kind';
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_array_textPartToInterpolated (textParts jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN coalesce((
        SELECT
            jsonb_agg(__migration_textPartToInterpolated (value))
        FROM jsonb_array_elements(textParts)), '[]');
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

UPDATE
    questions
SET
    enonce = (
        SELECT
            jsonb_agg(
                CASE WHEN (value -> 'Kind')::int = 11 THEN
                    -- OrderedListFieldBlockKind
                    jsonb_set(jsonb_set(value, '{Data, Answer}', __migration_array_textPartToInterpolated (value -> 'Data' -> 'Answer')), '{Data, AdditionalProposals}', __migration_array_textPartToInterpolated (value -> 'Data' -> 'AdditionalProposals'))
                ELSE
                    value
                END)
        FROM
            jsonb_array_elements(enonce));

DROP FUNCTION __migration_textPartToInterpolated;

DROP FUNCTION __migration_array_textPartToInterpolated;

-- Put the CONSTRAINT back
ALTER TABLE questions
    ADD CONSTRAINT enonce_structgen_validate_json_array_exe_Block CHECK (structgen_validate_json_array_exe_Block (enonce));

