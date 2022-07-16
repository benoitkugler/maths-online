-- v0.6.1
-- add Area entries to figures
-- blocks : FigureBlock (2), FigureAffineLineFieldBlock (1), FigurePointFieldBlock (3), FigureVectorFieldBlock (4), FigureVectorPairFieldBlock (5)

UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN (value -> 'Kind')::int = 2 THEN
                        jsonb_set(value, '{Data, Drawings, Areas}', '[]')
                    WHEN (value -> 'Kind')::int = 1 THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Areas}', '[]')
                    WHEN (value -> 'Kind')::int = 3 THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Areas}', '[]')
                    WHEN (value -> 'Kind')::int = 4 THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Areas}', '[]')
                    WHEN (value -> 'Kind')::int = 5 THEN
                        jsonb_set(value, '{Data, Figure, Drawings, Areas}', '[]')
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));

