-- v0.6.3
-- enable interpolation in ExpressionFieldBlock label

BEGIN;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN (value ->> 'Kind') = 'ExpressionFieldBlock' THEN
                        -- only text is actually used
                        jsonb_set(value, '{Data, Label}', value -> 'Data' -> 'Label' -> 'Content')
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));
COMMIT;

