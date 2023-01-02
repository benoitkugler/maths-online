-- v1.1.2
--
-- Replace randSymbol and choiceSymbol by randChoice and choiceFrom,
-- and update the syntax accordingly
-- randSymbol(A;B;C) -> randChoice(A;B;C)
-- choiceSymbol((A;B;C); k == 1) -> choiceFrom(A;B;C; k == 1)

CREATE OR REPLACE FUNCTION __migration_replace_randSymbol (expression text)
    RETURNS text
    AS $$
DECLARE
BEGIN
    -- randSymbol: just rename the function
    expression = REPLACE(expression, 'randSymbol', 'randChoice');
    expression = REPLACE(expression, 'randsymbol', 'randChoice');
    -- choiceSymbol : remove the parenthesis and rename
    expression = REGEXP_REPLACE(expression, 'choiceSymbol\(\((.+)\);', 'choiceFrom(\1; ', 'g');
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
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (__migration_replace_randSymbol (value ->> 'expression'))))
            FROM jsonb_array_elements(coalesce(page -> 'parameters' ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));

UPDATE
    exercices
SET
    parameters = jsonb_set(parameters, '{Variables}', coalesce((
            SELECT
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (__migration_replace_randSymbol (value ->> 'expression'))))
            FROM jsonb_array_elements(coalesce(parameters ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));

