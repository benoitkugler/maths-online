-- add IsDiscrete field to function graph
-- merge 4 geometric construction fields (also adding functions as background)

BEGIN;
--
CREATE OR REPLACE FUNCTION __migration_point_func (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN data || '{ "IsDiscrete" : false}'::jsonb;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_graphe_sequences (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN data || '{ "SequenceExprs" : []}'::jsonb;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_FigurePoint (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object(
        --
        'Field',
        --
        jsonb_build_object('Kind', 'GFPoint', 'Data', data - 'Figure'),
        --
        'Background',
        --
        jsonb_build_object('Kind', 'FigureBlock', 'Data', data -> 'Figure'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_FigureVector (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object(
        --
        'Field',
        --
        jsonb_build_object('Kind', 'GFVector', 'Data', data - 'Figure'),
        --
        'Background',
        --
        jsonb_build_object('Kind', 'FigureBlock', 'Data', data -> 'Figure'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_FigureAffineLine (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object(
        --
        'Field',
        --
        jsonb_build_object('Kind', 'GFAffineLine', 'Data', data - 'Figure'),
        --
        'Background',
        --
        jsonb_build_object('Kind', 'FigureBlock', 'Data', data -> 'Figure'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_FigureVectorPair (data jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_build_object(
        --
        'Field',
        --
        jsonb_build_object('Kind', 'GFVectorPair', 'Data', data - 'Figure'),
        --
        'Background',
        --
        jsonb_build_object('Kind', 'FigureBlock', 'Data', data -> 'Figure'));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
CREATE OR REPLACE FUNCTION __migration_block (value jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    CASE WHEN value ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
        RETURN jsonb_set(value, '{Data}', __migration_point_func (value -> 'Data'));
    WHEN value ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN jsonb_set(value, '{Data}', __migration_graphe_sequences (value -> 'Data'));
        -- merge geo fields
    WHEN value ->> 'Kind' = 'FigurePointFieldBlock' THEN
        RETURN jsonb_build_object('Kind', 'GeometricConstructionFieldBlock', 'Data', __migration_FigurePoint (value -> 'Data'));
    WHEN value ->> 'Kind' = 'FigureVectorFieldBlock' THEN
        RETURN jsonb_build_object('Kind', 'GeometricConstructionFieldBlock', 'Data', __migration_FigureVector (value -> 'Data'));
    WHEN value ->> 'Kind' = 'FigureAffineLineFieldBlock' THEN
        RETURN jsonb_build_object('Kind', 'GeometricConstructionFieldBlock', 'Data', __migration_FigureAffineLine (value -> 'Data'));
    WHEN value ->> 'Kind' = 'FigureVectorPairFieldBlock' THEN
        RETURN jsonb_build_object('Kind', 'GeometricConstructionFieldBlock', 'Data', __migration_FigureVectorPair (value -> 'Data'));
    ELSE
        RETURN value;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
--

UPDATE
    questions
SET
    enonce = coalesce((
        SELECT
            jsonb_agg(__migration_block (value))
        FROM jsonb_array_elements(enonce)), '[]'),
    correction = coalesce((
        SELECT
            jsonb_agg(__migration_block (value))
        FROM jsonb_array_elements(correction)), '[]');
--
--

DROP FUNCTION __migration_point_func (jsonb);
DROP FUNCTION __migration_graphe_sequences (jsonb);
DROP FUNCTION __migration_block (jsonb);
DROP FUNCTION __migration_FigurePoint (jsonb);
DROP FUNCTION __migration_FigureVector (jsonb);
DROP FUNCTION __migration_FigureAffineLine (jsonb);
DROP FUNCTION __migration_FigureVectorPair (jsonb);
COMMIT;

