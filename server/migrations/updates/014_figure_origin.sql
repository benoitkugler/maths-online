-- v0.6.1
-- add ShowOrigin field to figures
-- blocks : FigureBlock (2), FigureAffineLineFieldBlock (1), FigurePointFieldBlock (3), FigureVectorFieldBlock (4), FigureVectorPairFieldBlock (5)

BEGIN;
-- First remove the constraint (should be added back after the migration)
ALTER TABLE questions
    DROP CONSTRAINT page_structgen_validate_json_que_QuestionPage;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN (value -> 'Kind')::int = 2 THEN
                        jsonb_set(value, '{Data, ShowOrigin}', to_jsonb (TRUE))
                    WHEN (value -> 'Kind')::int = 1 THEN
                        jsonb_set(value, '{Data, Figure, ShowOrigin}', to_jsonb (TRUE))
                    WHEN (value -> 'Kind')::int = 3 THEN
                        jsonb_set(value, '{Data, Figure, ShowOrigin}', to_jsonb (TRUE))
                    WHEN (value -> 'Kind')::int = 4 THEN
                        jsonb_set(value, '{Data, Figure, ShowOrigin}', to_jsonb (TRUE))
                    WHEN (value -> 'Kind')::int = 5 THEN
                        jsonb_set(value, '{Data, Figure, ShowOrigin}', to_jsonb (TRUE))
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));
ALTER TABLE questions
    ADD CONSTRAINT page_structgen_validate_json_que_QuestionPage CHECK (structgen_validate_json_que_QuestionPage (page));
COMMIT;

