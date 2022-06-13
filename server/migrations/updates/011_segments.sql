-- v0.5.0
-- add color and kind field to segments, and color to points
-- blocks : FigureBlock (2), FigureAffineLineFieldBlock (1), FigurePointFieldBlock (3), FigureVectorFieldBlock (4), FigureVectorPairFieldBlock (5)

CREATE OR REPLACE FUNCTION __migration_point (labeledPoint jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_set(labeledPoint, '{Color}', to_jsonb ('#FF0000'::text));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_points (namedRandomLabeledPoints jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN coalesce((
        SELECT
            (jsonb_agg(jsonb_set(value, '{Point}', __migration_point (value -> 'Point'))))
        FROM jsonb_array_elements(namedRandomLabeledPoints)), '[]');
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_segmentKind (asVector jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    CASE WHEN (asVector::boolean) = TRUE THEN
        RETURN 1;
    ELSE
        RETURN 0;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_segment (segment jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_set(jsonb_set(segment - 'AsVector', '{Color}', to_jsonb ('#FF0000'::text)), '{Kind}', __migration_segmentKind (segment -> 'AsVector'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_segments (segments jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN coalesce((
        SELECT
            (jsonb_agg(__migration_segment (value)))
        FROM jsonb_array_elements(segments)), '[]');
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_randomDrawings (drawings jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_set(jsonb_set(drawings, '{Points}', __migration_points (drawings -> 'Points')), '{Segments}', __migration_segments (drawings -> 'Segments'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

BEGIN;
-- First remove the constraint (should be added back after the migration)
ALTER TABLE questions
    DROP CONSTRAINT question_structgen_validate_json_exe_QuestionPage;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', (
            SELECT
                jsonb_agg(
                    CASE WHEN (value -> 'Kind')::int = 2 THEN
                        jsonb_set(value, '{Data, Drawings}', __migration_randomDrawings (value -> 'Data' -> 'Drawings'))
                    WHEN (value -> 'Kind')::int = 1 THEN
                        jsonb_set(value, '{Data, Figure, Drawings}', __migration_randomDrawings (value -> 'Data' -> 'Figure' -> 'Drawings'))
                    WHEN (value -> 'Kind')::int = 3 THEN
                        jsonb_set(value, '{Data, Figure, Drawings}', __migration_randomDrawings (value -> 'Data' -> 'Figure' -> 'Drawings'))
                    WHEN (value -> 'Kind')::int = 4 THEN
                        jsonb_set(value, '{Data, Figure, Drawings}', __migration_randomDrawings (value -> 'Data' -> 'Figure' -> 'Drawings'))
                    WHEN (value -> 'Kind')::int = 5 THEN
                        jsonb_set(value, '{Data, Figure, Drawings}', __migration_randomDrawings (value -> 'Data' -> 'Figure' -> 'Drawings'))
                    ELSE
                        value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')));
-- Put the CONSTRAINT back, after updating the definitions
ALTER TABLE questions
    ADD CONSTRAINT page_structgen_validate_json_exe_QuestionPage CHECK (structgen_validate_json_exe_QuestionPage (page));
COMMIT;

