BEGIN;
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
    Index integer NOT NULL,
    IdQuestion integer NOT NULL
);
ALTER TABLE random_monoquestion_variants
    ADD UNIQUE (IdStudent, IdRandomMonoquestion, INDEX);
ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdStudent) REFERENCES students;
ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions;
ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;
ALTER TABLE random_monoquestions
    ADD CHECK (NbRepeat > 0);
ALTER TABLE random_monoquestions
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups;
ALTER TABLE tasks
    ADD COLUMN IdRandomMonoquestion integer;
ALTER TABLE tasks
    DROP CONSTRAINT "tasks_check";
ALTER TABLE tasks
    DROP CONSTRAINT "tasks_check1";
ALTER TABLE tasks
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions;
ALTER TABLE tasks
    ADD CHECK ((IdExercice IS NOT NULL)::int + (IdMonoquestion IS NOT NULL)::int + (IdRandomMonoquestion IS NOT NULL)::int = 1);
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
DROP TABLE progression_questions;
DROP TABLE progressions;
CREATE TABLE progressions (
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL,
    Index integer NOT NULL,
    History boolean[]
);
ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask, INDEX);
ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;
ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;
COMMIT;

