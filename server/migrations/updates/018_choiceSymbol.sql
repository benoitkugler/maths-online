-- v0.6.1
-- rename randLetter to randSymbol to better match the new choiceSymbol

UPDATE
    questions
SET
    page = jsonb_set(page, '{parameters, Variables}', coalesce((
            SELECT
                jsonb_agg(jsonb_set(value, '{expression}', to_jsonb (REPLACE(REPLACE(value ->> 'expression', 'randLetter', 'randSymbol'), 'randletter', 'randSymbol'))))
            FROM jsonb_array_elements(coalesce(page -> 'parameters' ->> 'Variables', '[]')::jsonb)), '[]'::jsonb));

