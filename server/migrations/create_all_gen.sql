
-- sql/teacher/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE classrooms (
    Id serial PRIMARY KEY,
    IdTeacher integer NOT NULL,
    Name text NOT NULL,
    MaxRankThreshold integer NOT NULL
);

CREATE TABLE classroom_codes (
    IdClassroom integer NOT NULL,
    Code text NOT NULL,
    ExpiresAt timestamp(0) with time zone NOT NULL
);

CREATE TABLE students (
    Id serial PRIMARY KEY,
    Name text NOT NULL,
    Surname text NOT NULL,
    Birthday timestamp(0) with time zone NOT NULL,
    IdClassroom integer NOT NULL,
    Clients jsonb NOT NULL
);

CREATE TABLE teachers (
    Id serial PRIMARY KEY,
    Mail text NOT NULL,
    PasswordCrypted bytea NOT NULL,
    IsAdmin boolean NOT NULL,
    HasSimplifiedEditor boolean NOT NULL,
    Contact jsonb NOT NULL,
    FavoriteMatiere text CHECK (FavoriteMatiere IN ('ALLEMAND', 'ANGLAIS', 'AUTRE', 'ESPAGNOL', 'FRANCAIS', 'HISTOIRE-GEO', 'ITALIEN', 'MATHS', 'PHYSIQUE', 'SES', 'SVT')) NOT NULL
);

-- constraints
ALTER TABLE teachers
    ADD UNIQUE (Mail);

ALTER TABLE classrooms
    ADD UNIQUE (Id, IdTeacher);

ALTER TABLE classrooms
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE classroom_codes
    ADD UNIQUE (Code);

ALTER TABLE classroom_codes
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE students
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

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

ALTER TABLE students
    ADD CONSTRAINT Clients_gomacro CHECK (gomacro_validate_json_array_teac_Client (Clients));

ALTER TABLE teachers
    ADD CONSTRAINT Contact_gomacro CHECK (gomacro_validate_json_teac_Contact (Contact));

-- sql/editor/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE exercices (
    Id serial PRIMARY KEY,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Parameters jsonb NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL
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
    IdExercicegroup integer NOT NULL,
    Section integer CHECK (Section IN (2, 1, 5, 4, 3)) NOT NULL
);

CREATE TABLE questions (
    Id serial PRIMARY KEY,
    Subtitle text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer,
    Enonce jsonb NOT NULL,
    Parameters jsonb NOT NULL,
    Correction jsonb NOT NULL
);

CREATE TABLE questiongroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE questiongroup_tags (
    Tag text NOT NULL,
    IdQuestiongroup integer NOT NULL,
    Section integer CHECK (Section IN (2, 1, 5, 4, 3)) NOT NULL
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
    ADD CHECK (Tag = upper(Tag));

CREATE UNIQUE INDEX QuestiongroupTag_level ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 1
    /* Section.Level */
;

CREATE UNIQUE INDEX QuestiongroupTag_chapter ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 2
    /* Section.Chapter */
;

CREATE UNIQUE INDEX QuestiongroupTag_matiere ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 5
    /* Section.Matiere */
;

ALTER TABLE questiongroup_tags
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups ON DELETE CASCADE;

ALTER TABLE exercicegroups
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE exercicegroup_tags
    ADD UNIQUE (IdExercicegroup, Tag);

ALTER TABLE exercicegroup_tags
    ADD CHECK (Tag = upper(Tag));

CREATE UNIQUE INDEX ExercicegroupTag_level ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 1
    /* Section.Level */
;

CREATE UNIQUE INDEX ExercicegroupTag_chapter ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 2
    /* Section.Chapter */
;

CREATE UNIQUE INDEX ExercicegroupTag_matiere ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 5
    /* Section.Matiere */
;

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

ALTER TABLE questions
    ADD CONSTRAINT Correction_gomacro CHECK (gomacro_validate_json_array_ques_Block (Correction));

ALTER TABLE questions
    ADD CONSTRAINT Enonce_gomacro CHECK (gomacro_validate_json_array_ques_Block (Enonce));

ALTER TABLE exercices
    ADD CONSTRAINT Parameters_gomacro CHECK (gomacro_validate_json_array_ques_ParameterEntry (Parameters));

ALTER TABLE questions
    ADD CONSTRAINT Parameters_gomacro CHECK (gomacro_validate_json_array_ques_ParameterEntry (Parameters));

-- sql/trivial/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE selfaccess_trivials (
    IdClassroom integer NOT NULL,
    IdTrivial integer NOT NULL,
    IdTeacher integer NOT NULL
);

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

ALTER TABLE selfaccess_trivials
    ADD FOREIGN KEY (IdClassroom, IdTeacher) REFERENCES Classrooms (Id, IdTeacher) ON DELETE CASCADE;

ALTER TABLE selfaccess_trivials
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE selfaccess_trivials
    ADD FOREIGN KEY (IdTrivial) REFERENCES trivials ON DELETE CASCADE;

ALTER TABLE selfaccess_trivials
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

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
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL,
    Index smallint NOT NULL,
    History boolean[]
);

CREATE TABLE random_monoquestions (
    Id serial PRIMARY KEY,
    IdQuestiongroup integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL,
    Difficulty jsonb NOT NULL
);

CREATE TABLE random_monoquestion_variants (
    IdStudent integer NOT NULL,
    IdRandomMonoquestion integer NOT NULL,
    Index smallint NOT NULL,
    IdQuestion integer NOT NULL
);

CREATE TABLE tasks (
    Id serial PRIMARY KEY,
    IdExercice integer,
    IdMonoquestion integer,
    IdRandomMonoquestion integer
);

-- constraints
ALTER TABLE monoquestions
    ADD CHECK (NbRepeat > 0);

ALTER TABLE monoquestions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE random_monoquestions
    ADD CHECK (NbRepeat > 0);

ALTER TABLE random_monoquestions
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups;

ALTER TABLE tasks
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE tasks
    ADD CHECK ((IdExercice IS NOT NULL)::int + (IdMonoquestion IS NOT NULL)::int + (IdRandomMonoquestion IS NOT NULL)::int = 1);

ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdMonoquestion) REFERENCES monoquestions;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions;

ALTER TABLE random_monoquestion_variants
    ADD UNIQUE (IdStudent, IdRandomMonoquestion, INDEX);

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions ON DELETE CASCADE;

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask, INDEX);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;

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

ALTER TABLE random_monoquestions
    ADD CONSTRAINT Difficulty_gomacro CHECK (gomacro_validate_json_array_edit_DifficultyTag (Difficulty));

-- sql/homework/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    IdTeacher integer NOT NULL,
    Level text NOT NULL,
    Anonymous integer,
    Public boolean NOT NULL,
    Matiere text CHECK (Matiere IN ('ALLEMAND', 'ANGLAIS', 'AUTRE', 'ESPAGNOL', 'FRANCAIS', 'HISTOIRE-GEO', 'ITALIEN', 'MATHS', 'PHYSIQUE', 'SES', 'SVT')) NOT NULL
);

CREATE TABLE sheet_tasks (
    IdSheet integer NOT NULL,
    Index integer NOT NULL,
    IdTask integer NOT NULL
);

CREATE TABLE travails (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    IdSheet integer NOT NULL,
    Noted boolean NOT NULL,
    Deadline timestamp(0) with time zone NOT NULL,
    ShowAfter timestamp(0) with time zone NOT NULL
);

CREATE TABLE travail_exceptions (
    IdStudent integer NOT NULL,
    IdTravail integer NOT NULL,
    Deadline timestamp(0) with time zone,
    IgnoreForMark boolean NOT NULL
);

-- constraints
ALTER TABLE travails
    ADD UNIQUE (Id, IdSheet);

ALTER TABLE travails
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE travails
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (Id, Anonymous) REFERENCES travails (IdSheet, Id) ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (Anonymous) REFERENCES travails ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD PRIMARY KEY (IdSheet, INDEX);

ALTER TABLE sheet_tasks
    ADD UNIQUE (IdTask);

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdTask) REFERENCES tasks;

ALTER TABLE travail_exceptions
    ADD UNIQUE (IdStudent, IdTravail);

ALTER TABLE travail_exceptions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE travail_exceptions
    ADD FOREIGN KEY (IdTravail) REFERENCES travails ON DELETE CASCADE;

-- sql/reviews/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE reviews (
    Id serial PRIMARY KEY,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);

CREATE TABLE review_exercices (
    IdReview integer NOT NULL,
    IdExercice integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);

CREATE TABLE review_participations (
    IdReview integer NOT NULL,
    IdTeacher integer NOT NULL,
    Approval integer CHECK (Approval IN (0, 1, 2)) NOT NULL,
    Comments jsonb NOT NULL
);

CREATE TABLE review_questions (
    IdReview integer NOT NULL,
    IdQuestion integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);

CREATE TABLE review_sheets (
    IdReview integer NOT NULL,
    IdSheet integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);

CREATE TABLE review_trivials (
    IdReview integer NOT NULL,
    IdTrivial integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2, 3)) NOT NULL
);

-- constraints
ALTER TABLE reviews
    ADD UNIQUE (Id, Kind);

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind) ON DELETE CASCADE;

ALTER TABLE review_questions
    ADD CHECK (Kind = 0
    /* ReviewKind.KQuestion */);

ALTER TABLE review_questions
    ADD UNIQUE (IdQuestion);

ALTER TABLE review_questions
    ADD UNIQUE (IdReview);

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questiongroups;

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind) ON DELETE CASCADE;

ALTER TABLE review_exercices
    ADD CHECK (Kind = 1
    /* ReviewKind.KExercice */);

ALTER TABLE review_exercices
    ADD UNIQUE (IdExercice);

ALTER TABLE review_exercices
    ADD UNIQUE (IdReview);

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdExercice) REFERENCES exercicegroups;

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind) ON DELETE CASCADE;

ALTER TABLE review_trivials
    ADD CHECK (Kind = 2
    /* ReviewKind.KTrivial */);

ALTER TABLE review_trivials
    ADD UNIQUE (IdTrivial);

ALTER TABLE review_trivials
    ADD UNIQUE (IdReview);

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdTrivial) REFERENCES trivials;

ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind) ON DELETE CASCADE;

ALTER TABLE review_sheets
    ADD CHECK (Kind = 3
    /* ReviewKind.KSheet */);

ALTER TABLE review_sheets
    ADD UNIQUE (IdSheet);

ALTER TABLE review_sheets
    ADD UNIQUE (IdReview);

ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_sheets
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets;

ALTER TABLE review_participations
    ADD UNIQUE (IdReview, IdTeacher);

ALTER TABLE review_participations
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_participations
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

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

ALTER TABLE review_participations
    ADD CONSTRAINT Comments_gomacro CHECK (gomacro_validate_json_array_revi_Comment (Comments));

-- sql/ceintures/gen_create.sql
-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE beltevolutions (
    IdStudent integer NOT NULL,
    Level integer CHECK (Level IN (0, 1, 2, 3)) NOT NULL,
    Advance jsonb NOT NULL,
    Stats jsonb NOT NULL
);

CREATE TABLE beltquestions (
    Id serial PRIMARY KEY,
    Domain integer CHECK (DOMAIN IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)) NOT NULL,
    Rank integer CHECK (Rank IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)) NOT NULL,
    Parameters jsonb NOT NULL,
    Enonce jsonb NOT NULL,
    Correction jsonb NOT NULL,
    Repeat integer NOT NULL,
    Title text NOT NULL
);

-- constraints
ALTER TABLE beltevolutions
    ADD UNIQUE (IdStudent);

ALTER TABLE beltevolutions
    ADD FOREIGN KEY (IdStudent) REFERENCES students;

ALTER TABLE beltquestions
    ADD CHECK (Repeat > 0);

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

ALTER TABLE beltevolutions
    ADD CONSTRAINT Stats_gomacro CHECK (gomacro_validate_json_array_12_array_11_cein_Stat (Stats));

ALTER TABLE beltevolutions
    ADD CONSTRAINT Advance_gomacro CHECK (gomacro_validate_json_array_12_cein_Rank (Advance));

ALTER TABLE beltquestions
    ADD CONSTRAINT Correction_gomacro CHECK (gomacro_validate_json_array_ques_Block (Correction));

ALTER TABLE beltquestions
    ADD CONSTRAINT Enonce_gomacro CHECK (gomacro_validate_json_array_ques_Block (Enonce));

ALTER TABLE beltquestions
    ADD CONSTRAINT Parameters_gomacro CHECK (gomacro_validate_json_array_ques_ParameterEntry (Parameters));

