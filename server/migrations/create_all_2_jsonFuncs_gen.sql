CREATE OR REPLACE FUNCTION gomacro_validate_json_array_teac_Client (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_teac_Client (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_string (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a string', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_teac_Client (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Device', 'Time'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Device')
        AND gomacro_validate_json_string (data -> 'Time');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_teac_Contact (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Name', 'URL'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Name')
        AND gomacro_validate_json_string (data -> 'URL');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_array_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_ques_TextPart (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_array_string (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_string (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_boolean (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_boolean (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_clie_SignSymbol (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_clie_SignSymbol (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_Block (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_Block (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionArea (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionArea (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionDefinition (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionDefinition (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionPoint (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionPoint (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionSign (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionSign (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_ParameterEntry (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_ParameterEntry (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_ProofAssertion (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_ProofAssertion (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_TextPart (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_TreeNodeAnswer (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_TreeNodeAnswer (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_VariationTableBlock (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_VariationTableBlock (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_NamedRandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_NamedRandomLabeledPoint (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomArea (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomArea (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomCircle (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomCircle (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomLine (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomLine (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomSegment (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomSegment (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_string (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_string (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_boolean (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'boolean';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a boolean', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_clie_Binary (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a clie_Binary', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_clie_SignSymbol (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a clie_SignSymbol', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_expr_Variable (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Indice', 'Name'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Indice')
        AND gomacro_validate_json_number (data -> 'Name');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_number (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a number', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_Block (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'ExpressionFieldBlock' THEN
        RETURN gomacro_validate_json_ques_ExpressionFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigureBlock' THEN
        RETURN gomacro_validate_json_ques_FigureBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FormulaBlock' THEN
        RETURN gomacro_validate_json_ques_FormulaBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionPointsFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionsGraphBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'GeometricConstructionFieldBlock' THEN
        RETURN gomacro_validate_json_ques_GeometricConstructionFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'NumberFieldBlock' THEN
        RETURN gomacro_validate_json_ques_NumberFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'OrderedListFieldBlock' THEN
        RETURN gomacro_validate_json_ques_OrderedListFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofFieldBlock' THEN
        RETURN gomacro_validate_json_ques_ProofFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'RadioFieldBlock' THEN
        RETURN gomacro_validate_json_ques_RadioFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SetFieldBlock' THEN
        RETURN gomacro_validate_json_ques_SetFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SignTableBlock' THEN
        RETURN gomacro_validate_json_ques_SignTableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SignTableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_SignTableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TableBlock' THEN
        RETURN gomacro_validate_json_ques_TableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_TableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TextBlock' THEN
        RETURN gomacro_validate_json_ques_TextBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TreeBlock' THEN
        RETURN gomacro_validate_json_ques_TreeBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TreeFieldBlock' THEN
        RETURN gomacro_validate_json_ques_TreeFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VariationTableBlock' THEN
        RETURN gomacro_validate_json_ques_VariationTableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VariationTableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_VariationTableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VectorFieldBlock' THEN
        RETURN gomacro_validate_json_ques_VectorFieldBlock (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ComparisonLevel (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (102, 2, 1, 0);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_ComparisonLevel', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_CoordExpression (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ExpressionFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Expression', 'Label', 'ComparisonLevel', 'ShowFractionHelp'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Expression')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_ques_ComparisonLevel (data -> 'ComparisonLevel')
        AND gomacro_validate_json_boolean (data -> 'ShowFractionHelp');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigureBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Drawings', 'Bounds', 'ShowGrid', 'ShowOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_repe_RandomDrawings (data -> 'Drawings')
        AND gomacro_validate_json_repe_RepereBounds (data -> 'Bounds')
        AND gomacro_validate_json_boolean (data -> 'ShowGrid')
        AND gomacro_validate_json_boolean (data -> 'ShowOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FiguresOrGraphs (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'FigureBlock' THEN
        RETURN gomacro_validate_json_ques_FigureBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionsGraphBlock (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FormulaBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Parts');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionArea (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Bottom', 'Top', 'Left', 'Right', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Bottom')
        AND gomacro_validate_json_string (data -> 'Top')
        AND gomacro_validate_json_string (data -> 'Left')
        AND gomacro_validate_json_string (data -> 'Right')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionDecoration (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionDefinition (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Function', 'Decoration', 'Variable', 'From', 'To'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_ques_FunctionDecoration (data -> 'Decoration')
        AND gomacro_validate_json_expr_Variable (data -> 'Variable')
        AND gomacro_validate_json_string (data -> 'From')
        AND gomacro_validate_json_string (data -> 'To');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Function', 'X', 'Color', 'Legend'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_string (data -> 'Legend');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionPointsFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('IsDiscrete', 'Function', 'Label', 'Variable', 'XGrid'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_boolean (data -> 'IsDiscrete')
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_expr_Variable (data -> 'Variable')
        AND gomacro_validate_json_array_string (data -> 'XGrid');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionSign (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'FxSymbols', 'Signs'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_clie_SignSymbol (data -> 'FxSymbols')
        AND gomacro_validate_json_array_boolean (data -> 'Signs');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionsGraphBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('FunctionExprs', 'FunctionVariations', 'SequenceExprs', 'Areas', 'Points'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'FunctionExprs')
        AND gomacro_validate_json_array_ques_VariationTableBlock (data -> 'FunctionVariations')
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'SequenceExprs')
        AND gomacro_validate_json_array_ques_FunctionArea (data -> 'Areas')
        AND gomacro_validate_json_array_ques_FunctionPoint (data -> 'Points');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFAffineLine (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'A', 'B'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'A')
        AND gomacro_validate_json_string (data -> 'B');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFVector (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AnswerOrigin', 'MustHaveOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_ques_CoordExpression (data -> 'AnswerOrigin')
        AND gomacro_validate_json_boolean (data -> 'MustHaveOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFVectorPair (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Criterion'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_VectorPairCriterion (data -> 'Criterion');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GeoField (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'GFAffineLine' THEN
        RETURN gomacro_validate_json_ques_GFAffineLine (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFPoint' THEN
        RETURN gomacro_validate_json_ques_GFPoint (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFVector' THEN
        RETURN gomacro_validate_json_ques_GFVector (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFVectorPair' THEN
        RETURN gomacro_validate_json_ques_GFVectorPair (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GeometricConstructionFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Field', 'Background'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_GeoField (data -> 'Field')
        AND gomacro_validate_json_ques_FiguresOrGraphs (data -> 'Background');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_NumberFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Expression'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Expression');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_OrderedListFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Answer', 'AdditionalProposals'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'AdditionalProposals');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ParameterEntry (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'Co' THEN
        RETURN gomacro_validate_json_string (data -> 'Data');
    WHEN data ->> 'Kind' = 'In' THEN
        RETURN gomacro_validate_json_string (data -> 'Data');
    WHEN data ->> 'Kind' = 'Rp' THEN
        RETURN gomacro_validate_json_ques_Rp (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofAssertion (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'ProofEquality' THEN
        RETURN gomacro_validate_json_ques_ProofEquality (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofInvalid' THEN
        RETURN gomacro_validate_json_ques_ProofInvalid (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofNode' THEN
        RETURN gomacro_validate_json_ques_ProofNode (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofSequence' THEN
        RETURN gomacro_validate_json_ques_ProofSequence (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofStatement' THEN
        RETURN gomacro_validate_json_ques_ProofStatement (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofEquality (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Terms'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Terms');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_ProofSequence (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofInvalid (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofNode (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Left', 'Right', 'Op'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_ProofAssertion (data -> 'Left')
        AND gomacro_validate_json_ques_ProofAssertion (data -> 'Right')
        AND gomacro_validate_json_clie_Binary (data -> 'Op');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofSequence (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_ProofAssertion (data -> 'Parts');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofStatement (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Content'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Content');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_RadioFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'Proposals', 'AsDropDown'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'Proposals')
        AND gomacro_validate_json_boolean (data -> 'AsDropDown');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_Rp (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('expression', 'variable'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'expression')
        AND gomacro_validate_json_expr_Variable (data -> 'variable');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SetFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AdditionalSets'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'AdditionalSets');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SignTableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Xs', 'Functions'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_string (data -> 'Xs')
        AND gomacro_validate_json_array_ques_FunctionSign (data -> 'Functions');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SignTableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_SignTableBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('HorizontalHeaders', 'VerticalHeaders', 'Values'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TextPart (data -> 'HorizontalHeaders')
        AND gomacro_validate_json_array_ques_TextPart (data -> 'VerticalHeaders')
        AND gomacro_validate_json_array_array_ques_TextPart (data -> 'Values');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('HorizontalHeaders', 'VerticalHeaders', 'Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TextPart (data -> 'HorizontalHeaders')
        AND gomacro_validate_json_array_ques_TextPart (data -> 'VerticalHeaders')
        AND gomacro_validate_json_array_array_string (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts', 'Bold', 'Italic', 'Smaller'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Parts')
        AND gomacro_validate_json_boolean (data -> 'Bold')
        AND gomacro_validate_json_boolean (data -> 'Italic')
        AND gomacro_validate_json_boolean (data -> 'Smaller');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextKind (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_TextKind', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Content', 'Kind'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Content')
        AND gomacro_validate_json_ques_TextKind (data -> 'Kind');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('EventsProposals', 'AnswerRoot'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_string (data -> 'EventsProposals')
        AND gomacro_validate_json_ques_TreeNodeAnswer (data -> 'AnswerRoot');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_TreeBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeNodeAnswer (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Children', 'Probabilities', 'Value'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TreeNodeAnswer (data -> 'Children')
        AND gomacro_validate_json_array_string (data -> 'Probabilities')
        AND gomacro_validate_json_number (data -> 'Value');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VariationTableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Xs', 'Fxs'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_string (data -> 'Xs')
        AND gomacro_validate_json_array_string (data -> 'Fxs');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VariationTableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_VariationTableBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VectorFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AcceptColinear', 'DisplayColumn'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_boolean (data -> 'AcceptColinear')
        AND gomacro_validate_json_boolean (data -> 'DisplayColumn');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VectorPairCriterion (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_VectorPairCriterion', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_Coord (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_number (data -> 'X')
        AND gomacro_validate_json_number (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_LabelPos (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2, 3, 4, 5, 6, 7, 8);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a repe_LabelPos', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_NamedRandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Name', 'Point'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Name')
        AND gomacro_validate_json_repe_RandomLabeledPoint (data -> 'Point');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomArea (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Color', 'Points'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_array_string (data -> 'Points');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomCircle (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Center', 'Radius', 'LineColor', 'FillColor', 'Legend'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_repe_RandomCoord (data -> 'Center')
        AND gomacro_validate_json_string (data -> 'Radius')
        AND gomacro_validate_json_string (data -> 'LineColor')
        AND gomacro_validate_json_string (data -> 'FillColor')
        AND gomacro_validate_json_string (data -> 'Legend');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomCoord (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomDrawings (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Points', 'Segments', 'Lines', 'Circles', 'Areas'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_repe_NamedRandomLabeledPoint (data -> 'Points')
        AND gomacro_validate_json_array_repe_RandomSegment (data -> 'Segments')
        AND gomacro_validate_json_array_repe_RandomLine (data -> 'Lines')
        AND gomacro_validate_json_array_repe_RandomCircle (data -> 'Circles')
        AND gomacro_validate_json_array_repe_RandomArea (data -> 'Areas');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Color', 'Coord', 'Pos'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_repe_RandomCoord (data -> 'Coord')
        AND gomacro_validate_json_repe_LabelPos (data -> 'Pos');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomLine (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'A', 'B', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'A')
        AND gomacro_validate_json_string (data -> 'B')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomSegment (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('LabelName', 'From', 'To', 'Color', 'LabelPos', 'Kind'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'LabelName')
        AND gomacro_validate_json_string (data -> 'From')
        AND gomacro_validate_json_string (data -> 'To')
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_repe_LabelPos (data -> 'LabelPos')
        AND gomacro_validate_json_repe_SegmentKind (data -> 'Kind');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RepereBounds (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Width', 'Height', 'Origin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_number (data -> 'Width')
        AND gomacro_validate_json_number (data -> 'Height')
        AND gomacro_validate_json_repe_Coord (data -> 'Origin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_SegmentKind (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a repe_SegmentKind', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_string (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a string', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_5_array_array_edit_TagSection (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_array_edit_TagSection (value))
        FROM
            jsonb_array_elements(data))
        AND jsonb_array_length(data) = 5;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_array_edit_TagSection (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_edit_TagSection (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_edit_DifficultyTag (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_edit_DifficultyTag (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_edit_TagSection (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_edit_TagSection (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_edit_DifficultyTag (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string'
    AND data #>> '{}' IN ('★', '★★', '★★★', '');
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a edit_DifficultyTag', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_edit_Section (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (2, 1, 5, 4, 3);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a edit_Section', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_edit_TagSection (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Tag', 'Section'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Tag')
        AND gomacro_validate_json_edit_Section (data -> 'Section');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_string (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a string', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_triv_CategoriesQuestions (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Tags', 'Difficulties'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_5_array_array_edit_TagSection (data -> 'Tags')
        AND gomacro_validate_json_array_edit_DifficultyTag (data -> 'Difficulties');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_edit_DifficultyTag (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_edit_DifficultyTag (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_edit_DifficultyTag (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string'
    AND data #>> '{}' IN ('★', '★★', '★★★', '');
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a edit_DifficultyTag', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_revi_Comment (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_revi_Comment (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_revi_Comment (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Time', 'Message'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Time')
        AND gomacro_validate_json_string (data -> 'Message');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_string (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a string', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_11_cein_Stat (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_cein_Stat (value))
        FROM
            jsonb_array_elements(data))
        AND jsonb_array_length(data) = 11;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_12_array_11_cein_Stat (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_11_cein_Stat (value))
        FROM
            jsonb_array_elements(data))
        AND jsonb_array_length(data) = 12;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_12_cein_Rank (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_cein_Rank (value))
        FROM
            jsonb_array_elements(data))
        AND jsonb_array_length(data) = 12;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_array_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_ques_TextPart (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_array_string (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_string (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_boolean (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_boolean (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_clie_SignSymbol (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_clie_SignSymbol (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_Block (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_Block (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionArea (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionArea (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionDefinition (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionDefinition (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionPoint (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionPoint (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_FunctionSign (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_FunctionSign (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_ParameterEntry (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_ParameterEntry (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_ProofAssertion (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_ProofAssertion (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_TextPart (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_TreeNodeAnswer (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_TreeNodeAnswer (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_VariationTableBlock (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_ques_VariationTableBlock (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_NamedRandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_NamedRandomLabeledPoint (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomArea (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomArea (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomCircle (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomCircle (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomLine (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomLine (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_repe_RandomSegment (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_repe_RandomSegment (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_string (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) = 'null' THEN
        RETURN TRUE;
    END IF;
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    IF jsonb_array_length(data) = 0 THEN
        RETURN TRUE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_string (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_boolean (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'boolean';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a boolean', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_cein_Rank (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a cein_Rank', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_cein_Stat (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Success', 'Failure'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_number (data -> 'Success')
        AND gomacro_validate_json_number (data -> 'Failure');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_clie_Binary (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a clie_Binary', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_clie_SignSymbol (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a clie_SignSymbol', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_expr_Variable (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Indice', 'Name'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Indice')
        AND gomacro_validate_json_number (data -> 'Name');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_number (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a number', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_Block (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'ExpressionFieldBlock' THEN
        RETURN gomacro_validate_json_ques_ExpressionFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigureBlock' THEN
        RETURN gomacro_validate_json_ques_FigureBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FormulaBlock' THEN
        RETURN gomacro_validate_json_ques_FormulaBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionPointsFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionsGraphBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'GeometricConstructionFieldBlock' THEN
        RETURN gomacro_validate_json_ques_GeometricConstructionFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'NumberFieldBlock' THEN
        RETURN gomacro_validate_json_ques_NumberFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'OrderedListFieldBlock' THEN
        RETURN gomacro_validate_json_ques_OrderedListFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofFieldBlock' THEN
        RETURN gomacro_validate_json_ques_ProofFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'RadioFieldBlock' THEN
        RETURN gomacro_validate_json_ques_RadioFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SetFieldBlock' THEN
        RETURN gomacro_validate_json_ques_SetFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SignTableBlock' THEN
        RETURN gomacro_validate_json_ques_SignTableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'SignTableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_SignTableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TableBlock' THEN
        RETURN gomacro_validate_json_ques_TableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_TableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TextBlock' THEN
        RETURN gomacro_validate_json_ques_TextBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TreeBlock' THEN
        RETURN gomacro_validate_json_ques_TreeBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'TreeFieldBlock' THEN
        RETURN gomacro_validate_json_ques_TreeFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VariationTableBlock' THEN
        RETURN gomacro_validate_json_ques_VariationTableBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VariationTableFieldBlock' THEN
        RETURN gomacro_validate_json_ques_VariationTableFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'VectorFieldBlock' THEN
        RETURN gomacro_validate_json_ques_VectorFieldBlock (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ComparisonLevel (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (102, 2, 1, 0);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_ComparisonLevel', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_CoordExpression (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ExpressionFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Expression', 'Label', 'ComparisonLevel', 'ShowFractionHelp'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Expression')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_ques_ComparisonLevel (data -> 'ComparisonLevel')
        AND gomacro_validate_json_boolean (data -> 'ShowFractionHelp');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigureBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Drawings', 'Bounds', 'ShowGrid', 'ShowOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_repe_RandomDrawings (data -> 'Drawings')
        AND gomacro_validate_json_repe_RepereBounds (data -> 'Bounds')
        AND gomacro_validate_json_boolean (data -> 'ShowGrid')
        AND gomacro_validate_json_boolean (data -> 'ShowOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FiguresOrGraphs (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'FigureBlock' THEN
        RETURN gomacro_validate_json_ques_FigureBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionsGraphBlock (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FormulaBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Parts');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionArea (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Bottom', 'Top', 'Left', 'Right', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Bottom')
        AND gomacro_validate_json_string (data -> 'Top')
        AND gomacro_validate_json_string (data -> 'Left')
        AND gomacro_validate_json_string (data -> 'Right')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionDecoration (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionDefinition (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Function', 'Decoration', 'Variable', 'From', 'To'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_ques_FunctionDecoration (data -> 'Decoration')
        AND gomacro_validate_json_expr_Variable (data -> 'Variable')
        AND gomacro_validate_json_string (data -> 'From')
        AND gomacro_validate_json_string (data -> 'To');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Function', 'X', 'Color', 'Legend'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_string (data -> 'Legend');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionPointsFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('IsDiscrete', 'Function', 'Label', 'Variable', 'XGrid'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_boolean (data -> 'IsDiscrete')
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_expr_Variable (data -> 'Variable')
        AND gomacro_validate_json_array_string (data -> 'XGrid');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionSign (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'FxSymbols', 'Signs'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_clie_SignSymbol (data -> 'FxSymbols')
        AND gomacro_validate_json_array_boolean (data -> 'Signs');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FunctionsGraphBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('FunctionExprs', 'FunctionVariations', 'SequenceExprs', 'Areas', 'Points'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'FunctionExprs')
        AND gomacro_validate_json_array_ques_VariationTableBlock (data -> 'FunctionVariations')
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'SequenceExprs')
        AND gomacro_validate_json_array_ques_FunctionArea (data -> 'Areas')
        AND gomacro_validate_json_array_ques_FunctionPoint (data -> 'Points');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFAffineLine (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'A', 'B'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'A')
        AND gomacro_validate_json_string (data -> 'B');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFVector (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AnswerOrigin', 'MustHaveOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_ques_CoordExpression (data -> 'AnswerOrigin')
        AND gomacro_validate_json_boolean (data -> 'MustHaveOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GFVectorPair (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Criterion'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_VectorPairCriterion (data -> 'Criterion');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GeoField (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'GFAffineLine' THEN
        RETURN gomacro_validate_json_ques_GFAffineLine (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFPoint' THEN
        RETURN gomacro_validate_json_ques_GFPoint (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFVector' THEN
        RETURN gomacro_validate_json_ques_GFVector (data -> 'Data');
    WHEN data ->> 'Kind' = 'GFVectorPair' THEN
        RETURN gomacro_validate_json_ques_GFVectorPair (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_GeometricConstructionFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Field', 'Background'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_GeoField (data -> 'Field')
        AND gomacro_validate_json_ques_FiguresOrGraphs (data -> 'Background');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_NumberFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Expression'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Expression');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_OrderedListFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Answer', 'AdditionalProposals'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'AdditionalProposals');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ParameterEntry (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'Co' THEN
        RETURN gomacro_validate_json_string (data -> 'Data');
    WHEN data ->> 'Kind' = 'In' THEN
        RETURN gomacro_validate_json_string (data -> 'Data');
    WHEN data ->> 'Kind' = 'Rp' THEN
        RETURN gomacro_validate_json_ques_Rp (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofAssertion (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'object' OR jsonb_typeof(data -> 'Kind') != 'string' OR jsonb_typeof(data -> 'Data') = 'null' THEN
        RETURN FALSE;
    END IF;
    CASE WHEN data ->> 'Kind' = 'ProofEquality' THEN
        RETURN gomacro_validate_json_ques_ProofEquality (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofInvalid' THEN
        RETURN gomacro_validate_json_ques_ProofInvalid (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofNode' THEN
        RETURN gomacro_validate_json_ques_ProofNode (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofSequence' THEN
        RETURN gomacro_validate_json_ques_ProofSequence (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofStatement' THEN
        RETURN gomacro_validate_json_ques_ProofStatement (data -> 'Data');
    ELSE
        RETURN FALSE;
    END CASE;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofEquality (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Terms'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Terms');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_ProofSequence (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofInvalid (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(TRUE)
        FROM
            jsonb_each(data));
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofNode (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Left', 'Right', 'Op'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_ProofAssertion (data -> 'Left')
        AND gomacro_validate_json_ques_ProofAssertion (data -> 'Right')
        AND gomacro_validate_json_clie_Binary (data -> 'Op');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofSequence (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_ProofAssertion (data -> 'Parts');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_ProofStatement (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Content'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Content');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_RadioFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'Proposals', 'AsDropDown'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'Proposals')
        AND gomacro_validate_json_boolean (data -> 'AsDropDown');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_Rp (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('expression', 'variable'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'expression')
        AND gomacro_validate_json_expr_Variable (data -> 'variable');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SetFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AdditionalSets'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Answer')
        AND gomacro_validate_json_array_string (data -> 'AdditionalSets');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SignTableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Xs', 'Functions'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_string (data -> 'Xs')
        AND gomacro_validate_json_array_ques_FunctionSign (data -> 'Functions');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SignTableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_SignTableBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('HorizontalHeaders', 'VerticalHeaders', 'Values'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TextPart (data -> 'HorizontalHeaders')
        AND gomacro_validate_json_array_ques_TextPart (data -> 'VerticalHeaders')
        AND gomacro_validate_json_array_array_ques_TextPart (data -> 'Values');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('HorizontalHeaders', 'VerticalHeaders', 'Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TextPart (data -> 'HorizontalHeaders')
        AND gomacro_validate_json_array_ques_TextPart (data -> 'VerticalHeaders')
        AND gomacro_validate_json_array_array_string (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Parts', 'Bold', 'Italic', 'Smaller'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Parts')
        AND gomacro_validate_json_boolean (data -> 'Bold')
        AND gomacro_validate_json_boolean (data -> 'Italic')
        AND gomacro_validate_json_boolean (data -> 'Smaller');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextKind (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_TextKind', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TextPart (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Content', 'Kind'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Content')
        AND gomacro_validate_json_ques_TextKind (data -> 'Kind');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('EventsProposals', 'AnswerRoot'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_string (data -> 'EventsProposals')
        AND gomacro_validate_json_ques_TreeNodeAnswer (data -> 'AnswerRoot');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_TreeBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_TreeNodeAnswer (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Children', 'Probabilities', 'Value'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_TreeNodeAnswer (data -> 'Children')
        AND gomacro_validate_json_array_string (data -> 'Probabilities')
        AND gomacro_validate_json_number (data -> 'Value');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VariationTableBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'Xs', 'Fxs'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_string (data -> 'Xs')
        AND gomacro_validate_json_array_string (data -> 'Fxs');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VariationTableFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_VariationTableBlock (data -> 'Answer');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VectorFieldBlock (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Answer', 'AcceptColinear', 'DisplayColumn'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_boolean (data -> 'AcceptColinear')
        AND gomacro_validate_json_boolean (data -> 'DisplayColumn');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_VectorPairCriterion (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_VectorPairCriterion', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_Coord (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_number (data -> 'X')
        AND gomacro_validate_json_number (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_LabelPos (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2, 3, 4, 5, 6, 7, 8);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a repe_LabelPos', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_NamedRandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Name', 'Point'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Name')
        AND gomacro_validate_json_repe_RandomLabeledPoint (data -> 'Point');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomArea (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Color', 'Points'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_array_string (data -> 'Points');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomCircle (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Center', 'Radius', 'LineColor', 'FillColor', 'Legend'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_repe_RandomCoord (data -> 'Center')
        AND gomacro_validate_json_string (data -> 'Radius')
        AND gomacro_validate_json_string (data -> 'LineColor')
        AND gomacro_validate_json_string (data -> 'FillColor')
        AND gomacro_validate_json_string (data -> 'Legend');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomCoord (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('X', 'Y'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'X')
        AND gomacro_validate_json_string (data -> 'Y');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomDrawings (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Points', 'Segments', 'Lines', 'Circles', 'Areas'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_repe_NamedRandomLabeledPoint (data -> 'Points')
        AND gomacro_validate_json_array_repe_RandomSegment (data -> 'Segments')
        AND gomacro_validate_json_array_repe_RandomLine (data -> 'Lines')
        AND gomacro_validate_json_array_repe_RandomCircle (data -> 'Circles')
        AND gomacro_validate_json_array_repe_RandomArea (data -> 'Areas');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomLabeledPoint (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Color', 'Coord', 'Pos'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_repe_RandomCoord (data -> 'Coord')
        AND gomacro_validate_json_repe_LabelPos (data -> 'Pos');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomLine (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Label', 'A', 'B', 'Color'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'A')
        AND gomacro_validate_json_string (data -> 'B')
        AND gomacro_validate_json_string (data -> 'Color');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RandomSegment (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('LabelName', 'From', 'To', 'Color', 'LabelPos', 'Kind'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'LabelName')
        AND gomacro_validate_json_string (data -> 'From')
        AND gomacro_validate_json_string (data -> 'To')
        AND gomacro_validate_json_string (data -> 'Color')
        AND gomacro_validate_json_repe_LabelPos (data -> 'LabelPos')
        AND gomacro_validate_json_repe_SegmentKind (data -> 'Kind');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_RepereBounds (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(key IN ('Width', 'Height', 'Origin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_number (data -> 'Width')
        AND gomacro_validate_json_number (data -> 'Height')
        AND gomacro_validate_json_repe_Coord (data -> 'Origin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_repe_SegmentKind (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a repe_SegmentKind', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_string (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'string';
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a string', data;
    END IF;
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

