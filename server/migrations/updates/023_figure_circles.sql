-- v0.6.3
-- add circle entries to figure

BEGIN;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN value ->> 'Kind' = 'FigureBlock' THEN
                        jsonb_set(value, '{Data, Drawings, Circles}', '[]')
                    WHEN value ->> 'Kind' = 'FigureAffineLineFieldBlock' THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Circles}', '[]')
                    WHEN value ->> 'Kind' = 'FigurePointFieldBlock' THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Circles}', '[]')
                    WHEN value ->> 'Kind' = 'FigureVectorFieldBlock' THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Circles}', '[]')
                    WHEN value ->> 'Kind' = 'FigureVectorPairFieldBlock' THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Circles}', '[]')
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));
COMMIT;

