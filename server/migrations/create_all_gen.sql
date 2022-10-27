
-- sql/teacher/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE classrooms (
    Id serial PRIMARY KEY,
    IdTeacher integer NOT NULL,
    Name text NOT NULL
);

CREATE TABLE students (
    Id serial PRIMARY KEY,
    Name text NOT NULL,
    Surname text NOT NULL,
    Birthday date NOT NULL,
    TrivialSuccess integer NOT NULL,
    IsClientAttached boolean NOT NULL,
    IdClassroom integer NOT NULL
);

CREATE TABLE teachers (
    Id serial PRIMARY KEY,
    Mail text NOT NULL,
    PasswordCrypted bytea NOT NULL,
    IsAdmin boolean NOT NULL
);

-- constraints
ALTER TABLE teachers
    ADD UNIQUE (Mail);

ALTER TABLE classrooms
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE students
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

-- sql/editor/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE exercices (
    Id serial PRIMARY KEY,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Parameters jsonb NOT NULL
);

CREATE TABLE exercice_questions (
    IdExercice integer NOT NULL,
    IdQuestion integer NOT NULL,
    Bareme integer NOT NULL,
    Index integer NOT NULL
);

CREATE TABLE exercicegroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE exercicegroup_tags (
    Tag text NOT NULL,
    IdExercicegroup integer NOT NULL
);

CREATE TABLE questions (
    Id serial PRIMARY KEY,
    Page jsonb NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer
);

CREATE TABLE questiongroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE questiongroup_tags (
    Tag text NOT NULL,
    IdQuestiongroup integer NOT NULL
);

-- constraints
ALTER TABLE questions
    ADD CHECK (NeedExercice IS NOT NULL
        OR IdGroup IS NOT NULL);

ALTER TABLE questions
    ADD UNIQUE (Id, NeedExercice);

ALTER TABLE questions
    ADD FOREIGN KEY (NeedExercice) REFERENCES exercices;

ALTER TABLE questions
    ADD FOREIGN KEY (IdGroup) REFERENCES questiongroups ON DELETE CASCADE;

ALTER TABLE questiongroups
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE questiongroup_tags
    ADD UNIQUE (IdQuestiongroup, Tag);

ALTER TABLE questiongroup_tags
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups ON DELETE CASCADE;

ALTER TABLE exercicegroups
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE exercicegroup_tags
    ADD UNIQUE (IdExercicegroup, Tag);

ALTER TABLE exercicegroup_tags
    ADD FOREIGN KEY (IdExercicegroup) REFERENCES exercicegroups ON DELETE CASCADE;

ALTER TABLE exercices
    ADD FOREIGN KEY (IdGroup) REFERENCES exercicegroups;

ALTER TABLE exercice_questions
    ADD PRIMARY KEY (IdExercice, INDEX);

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdExercice, IdQuestion) REFERENCES Questions (NeedExercice, Id);

ALTER TABLE exercice_questions
    ADD UNIQUE (IdQuestion);

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices ON DELETE CASCADE;

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

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

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_RandomParameter (data jsonb)
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
            bool_and(gomacro_validate_json_ques_RandomParameter (value))
        FROM
            jsonb_array_elements(data));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_ques_SignSymbol (data jsonb)
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
            bool_and(gomacro_validate_json_ques_SignSymbol (value))
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_func_FunctionDecoration (data jsonb)
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
    WHEN data ->> 'Kind' = 'FigureAffineLineFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FigureAffineLineFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigureBlock' THEN
        RETURN gomacro_validate_json_ques_FigureBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigurePointFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FigurePointFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigureVectorFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FigureVectorFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FigureVectorPairFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FigureVectorPairFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FormulaBlock' THEN
        RETURN gomacro_validate_json_ques_FormulaBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionPointsFieldBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionPointsFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'FunctionsGraphBlock' THEN
        RETURN gomacro_validate_json_ques_FunctionsGraphBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'NumberFieldBlock' THEN
        RETURN gomacro_validate_json_ques_NumberFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'OrderedListFieldBlock' THEN
        RETURN gomacro_validate_json_ques_OrderedListFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'ProofFieldBlock' THEN
        RETURN gomacro_validate_json_ques_ProofFieldBlock (data -> 'Data');
    WHEN data ->> 'Kind' = 'RadioFieldBlock' THEN
        RETURN gomacro_validate_json_ques_RadioFieldBlock (data -> 'Data');
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
            bool_and(key IN ('Expression', 'Label', 'ComparisonLevel'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Expression')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_ques_ComparisonLevel (data -> 'ComparisonLevel');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigureAffineLineFieldBlock (data jsonb)
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
            bool_and(key IN ('Label', 'A', 'B', 'Figure'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_string (data -> 'A')
        AND gomacro_validate_json_string (data -> 'B')
        AND gomacro_validate_json_ques_FigureBlock (data -> 'Figure');
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigurePointFieldBlock (data jsonb)
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
            bool_and(key IN ('Answer', 'Figure'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_ques_FigureBlock (data -> 'Figure');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigureVectorFieldBlock (data jsonb)
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
            bool_and(key IN ('Answer', 'AnswerOrigin', 'Figure', 'MustHaveOrigin'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_CoordExpression (data -> 'Answer')
        AND gomacro_validate_json_ques_CoordExpression (data -> 'AnswerOrigin')
        AND gomacro_validate_json_ques_FigureBlock (data -> 'Figure')
        AND gomacro_validate_json_boolean (data -> 'MustHaveOrigin');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_FigureVectorPairFieldBlock (data jsonb)
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
            bool_and(key IN ('Figure', 'Criterion'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_ques_FigureBlock (data -> 'Figure')
        AND gomacro_validate_json_ques_VectorPairCriterion (data -> 'Criterion');
    RETURN is_valid;
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
        AND gomacro_validate_json_func_FunctionDecoration (data -> 'Decoration')
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
            bool_and(key IN ('Function', 'Label', 'Variable', 'XGrid'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Function')
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_expr_Variable (data -> 'Variable')
        AND gomacro_validate_json_array_string (data -> 'XGrid');
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
            bool_and(key IN ('FunctionExprs', 'FunctionVariations', 'Areas', 'Points'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_FunctionDefinition (data -> 'FunctionExprs')
        AND gomacro_validate_json_array_ques_VariationTableBlock (data -> 'FunctionVariations')
        AND gomacro_validate_json_array_ques_FunctionArea (data -> 'Areas')
        AND gomacro_validate_json_array_ques_FunctionPoint (data -> 'Points');
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_Parameters (data jsonb)
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
            bool_and(key IN ('Variables', 'Intrinsics'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_RandomParameter (data -> 'Variables')
        AND gomacro_validate_json_array_string (data -> 'Intrinsics');
    RETURN is_valid;
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_QuestionPage (data jsonb)
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
            bool_and(key IN ('enonce', 'parameters'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_array_ques_Block (data -> 'enonce')
        AND gomacro_validate_json_ques_Parameters (data -> 'parameters');
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_RandomParameter (data jsonb)
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

CREATE OR REPLACE FUNCTION gomacro_validate_json_ques_SignSymbol (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean := jsonb_typeof(data) = 'number'
    AND data::int IN (0, 1, 2);
BEGIN
    IF NOT is_valid THEN
        RAISE WARNING '% is not a ques_SignSymbol', data;
    END IF;
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
            bool_and(key IN ('Label', 'FxSymbols', 'Xs', 'Signs'))
        FROM
            jsonb_each(data))
        AND gomacro_validate_json_string (data -> 'Label')
        AND gomacro_validate_json_array_ques_SignSymbol (data -> 'FxSymbols')
        AND gomacro_validate_json_array_string (data -> 'Xs')
        AND gomacro_validate_json_array_boolean (data -> 'Signs');
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

ALTER TABLE exercices
    ADD CONSTRAINT Parameters_gomacro CHECK (gomacro_validate_json_ques_Parameters (Parameters));

ALTER TABLE questions
    ADD CONSTRAINT Page_gomacro CHECK (gomacro_validate_json_ques_QuestionPage (Page));

-- sql/trivial/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE trivials (
    Id serial PRIMARY KEY,
    Questions jsonb NOT NULL,
    QuestionTimeout integer NOT NULL,
    ShowDecrassage boolean NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL,
    Name text NOT NULL
);

-- constraints
ALTER TABLE trivials
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

CREATE OR REPLACE FUNCTION gomacro_validate_json_array_5_array_array_string (data jsonb)
    RETURNS boolean
    AS $$
BEGIN
    IF jsonb_typeof(data) != 'array' THEN
        RETURN FALSE;
    END IF;
    RETURN (
        SELECT
            bool_and(gomacro_validate_json_array_array_string (value))
        FROM
            jsonb_array_elements(data))
        AND jsonb_array_length(data) = 5;
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
        AND gomacro_validate_json_array_5_array_array_string (data -> 'Tags')
        AND gomacro_validate_json_array_edit_DifficultyTag (data -> 'Difficulties');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

ALTER TABLE trivials
    ADD CONSTRAINT Questions_gomacro CHECK (gomacro_validate_json_triv_CategoriesQuestions (Questions));

-- sql/tasks/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE monoquestions (
    Id serial PRIMARY KEY,
    IdQuestion integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL
);

CREATE TABLE progressions (
    Id serial PRIMARY KEY,
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL
);

CREATE TABLE progression_questions (
    IdProgression integer NOT NULL,
    Index integer NOT NULL,
    History boolean[]
);

CREATE TABLE tasks (
    Id serial PRIMARY KEY,
    IdExercice integer,
    IdMonoquestion integer
);

-- constraints
ALTER TABLE monoquestions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE tasks
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE tasks
    ADD CHECK (IdExercice IS NOT NULL
        OR IdMonoquestion IS NOT NULL);

ALTER TABLE tasks
    ADD CHECK (IdExercice IS NULL
        OR IdMonoquestion IS NULL);

ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdMonoquestion) REFERENCES monoquestions;

ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdProgression) REFERENCES progressions ON DELETE CASCADE;

-- sql/homework/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    Title text NOT NULL,
    Notation integer CHECK (Notation IN (0, 1)) NOT NULL,
    Activated boolean NOT NULL,
    Deadline timestamp(0) with time zone NOT NULL
);

CREATE TABLE sheet_tasks (
    IdSheet integer NOT NULL,
    Index integer NOT NULL,
    IdTask integer NOT NULL
);

-- constraints
ALTER TABLE sheets
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD PRIMARY KEY (IdSheet, INDEX);

ALTER TABLE sheet_tasks
    ADD UNIQUE (IdTask);

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdTask) REFERENCES tasks;

