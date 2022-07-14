-- v0.6.1
-- replace int tags by string tags

UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN (value -> 'Kind')::int = 0 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('ExpressionFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 1 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FigureAffineLineFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 2 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FigureBlock'::text))
                    WHEN (value -> 'Kind')::int = 3 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FigurePointFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 4 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FigureVectorFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 5 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FigureVectorPairFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 6 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FormulaBlock'::text))
                    WHEN (value -> 'Kind')::int = 7 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FunctionGraphBlock'::text))
                    WHEN (value -> 'Kind')::int = 8 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FunctionPointsFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 9 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('FunctionVariationGraphBlock'::text))
                    WHEN (value -> 'Kind')::int = 10 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('NumberFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 11 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('OrderedListFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 12 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('RadioFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 13 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('SignTableBlock'::text))
                    WHEN (value -> 'Kind')::int = 14 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('TableBlock'::text))
                    WHEN (value -> 'Kind')::int = 15 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('TableFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 16 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('TextBlock'::text))
                    WHEN (value -> 'Kind')::int = 17 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('TreeFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 18 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('VariationTableBlock'::text))
                    WHEN (value -> 'Kind')::int = 19 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('VariationTableFieldBlock'::text))
                    WHEN (value -> 'Kind')::int = 20 THEN
                        jsonb_set(value, '{Kind}', to_jsonb ('VectorFieldBlock'::text))
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));

