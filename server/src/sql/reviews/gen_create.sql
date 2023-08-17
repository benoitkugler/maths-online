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

