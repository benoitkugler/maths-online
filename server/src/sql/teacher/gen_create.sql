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

