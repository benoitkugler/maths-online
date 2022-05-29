-- v0.4.0
-- move question content on its own type
-- the order of fields matter
-- requires the teacher table with one admin on it

ALTER TABLE questions
    ADD COLUMN page jsonb;

UPDATE
    questions
SET
    page = jsonb_build_object('title', to_jsonb (title), 'enonce', enonce, 'parameters', parameters);

ALTER TABLE questions
    ALTER COLUMN page SET NOT NULL;

CREATE OR REPLACE FUNCTION structgen_validate_json_exe_QuestionPage (data jsonb)
    RETURNS boolean
    AS $$
DECLARE
    is_valid boolean;
BEGIN
    IF jsonb_typeof(data) != 'object' THEN
        RETURN FALSE;
    END IF;
    is_valid := (
        SELECT
            bool_and(KEY IN ('title', 'enonce', 'parameters'))
        FROM
            jsonb_each(data))
        AND structgen_validate_json_string (data -> 'title')
        AND structgen_validate_json_array_exe_Block (data -> 'enonce')
        AND structgen_validate_json_exe_Parameters (data -> 'parameters');
    RETURN is_valid;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

ALTER TABLE questions
    ADD CONSTRAINT question_structgen_validate_json_exe_QuestionPage CHECK (structgen_validate_json_exe_QuestionPage (page));

ALTER TABLE questions
    DROP COLUMN title;

ALTER TABLE questions
    DROP COLUMN enonce;

ALTER TABLE questions
    DROP COLUMN parameters;

-- visibility
ALTER TABLE questions
    ADD COLUMN public boolean;

UPDATE
    questions
SET
    public = TRUE;

ALTER TABLE questions
    ALTER COLUMN public SET NOT NULL;

-- admin : all the existing questions are owned by the admin account
ALTER TABLE questions
    ADD COLUMN id_teacher integer;

UPDATE
    questions
SET
    id_teacher = 1;

ALTER TABLE questions
    ALTER COLUMN id_teacher SET NOT NULL;

ALTER TABLE questions
    ADD FOREIGN KEY (id_teacher) REFERENCES teachers;

