-- v1.2.2
-- replace @_<latex> by "<latex>"
-- check that @ is not used in questions

SELECT DISTINCT
    position('@' IN page ->> 'enonce')
FROM
    questions;

BEGIN;
--
-- implement the transformation for one expression
--

CREATE OR REPLACE FUNCTION __migration_replace_at_ (expression text)
    RETURNS text
    AS $$
DECLARE
BEGIN
    -- replace @_ by quotes
    expression = REGEXP_REPLACE(expression, '@_([A-z0-9\\{}]*)', '"\1"', 'g');
    -- IF position('"' IN expression) > 0 THEN
    --     RAISE WARNING '%', expression;
    -- END IF;
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
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (__migration_replace_at_ (value ->> 'expression'))))
            FROM jsonb_array_elements(coalesce(page -> 'parameters' ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));
UPDATE
    exercices
SET
    parameters = jsonb_set(parameters, '{Variables}', coalesce((
            SELECT
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (__migration_replace_at_ (value ->> 'expression'))))
            FROM jsonb_array_elements(coalesce(parameters ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));
DROP FUNCTION __migration_replace_at_ (expression text);
COMMIT;

