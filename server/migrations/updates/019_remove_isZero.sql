--v0.6.1
-- remove isZero in favor of the more natural == operator

CREATE OR REPLACE FUNCTION __migration_replace_isZero (expression text)
    RETURNS text
    AS $$
DECLARE
BEGIN
    -- note that the order matter
    expression = REGEXP_REPLACE(expression, '^isZero\((.+)-\s*(\d+)\)$', '\1 == \2', 'g');
    expression = REGEXP_REPLACE(expression, '^isZero\((.+)\+\s*(\d+)\)$', '\1 == -\2', 'g');
    expression = REGEXP_REPLACE(expression, '^isZero\((.)\)$', '\1 == 0', 'g');
    expression = REGEXP_REPLACE(expression, '^isZero\((.+)\)$', '(\1) == 0', 'g');
    expression = REGEXP_REPLACE(expression, 'isZero\((.+)\)', '((\1) == 0)', 'g');
    RETURN expression;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

UPDATE
    questions
SET
    page = jsonb_set(page, '{parameters, Variables}', coalesce((
            SELECT
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (__migration_replace_isZero (value ->> 'expression'))))
            FROM jsonb_array_elements(coalesce(page -> 'parameters' ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));

DROP FUNCTION __migration_replace_isZero (text);

