--v0.9.0
-- add points to FunctionsGraph
-- the constraints must be restore after the update scripts

ALTER TABLE questions
    DROP CONSTRAINT Page_gomacro;

BEGIN;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN value ->> 'Kind' = 'FunctionsGraphBlock' THEN
                        jsonb_set(value, '{Data, Points}', '[]')
                ELSE
                    value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));
COMMIT;

