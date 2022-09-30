--v0.8.0
-- add explicit question groups, add exercices groups, add monoquestions
---
-- constraints must be added after this script
--

BEGIN;
--
-- update question table
--

CREATE TABLE questiongroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);
-- compute implicit groups with the title.
-- To handle groups with mixed public attribute, we compute it later,
-- defaulting to False

INSERT INTO questiongroups (title, public, idteacher)
SELECT DISTINCT
    -- select standalone questions
    page ->> 'title',
    FALSE,
    idteacher
FROM
    questions
WHERE
    NeedExercice IS NULL;
-- publish group if all the questions are public
UPDATE
    questiongroups
SET
    public = TRUE
WHERE (
    SELECT
        bool_and(questions.public)
    FROM
        questions
    WHERE
        questiongroups.title = questions.page ->> 'title'
        AND questiongroups.idteacher = questions.idteacher);
-- attach questions to groups, by creating a new table, required to insert difficulty at the correct place
CREATE TABLE questions2 (
    Id serial,
    Page jsonb NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer
);
INSERT INTO questions2
SELECT
    id,
    Page,
    Page ->> 'title', -- Subtitle, modified after
    Description,
    '', --diffculty : completed after
    NeedExercice,
    (
        CASE WHEN NeedExercice IS NULL THEN
        ( SELECT DISTINCT
                questiongroups.id
            FROM
                questiongroups
            WHERE
                questiongroups.title = questions.page ->> 'title'
                AND questiongroups.idteacher = questions.idteacher)
        ELSE
            NULL
        END)
FROM
    questions;
--
-- extract difficulty
--

UPDATE
    questions2
SET
    Difficulty = '★'
WHERE
    id = ANY (
        SELECT
            idquestion
        FROM
            question_tags
        WHERE
            tag = '★');
UPDATE
    questions2
SET
    Difficulty = '★★'
WHERE
    id = ANY (
        SELECT
            idquestion
        FROM
            question_tags
        WHERE
            tag = '★★');
UPDATE
    questions2
SET
    Difficulty = '★★★'
WHERE
    id = ANY (
        SELECT
            idquestion
        FROM
            question_tags
        WHERE
            tag = '★★★');
--
-- split the tags between common tags, and subtitles.
-- also, remove the now useless difficulty tags
--

DELETE FROM question_tags
WHERE tag IN ('★', '★★', '★★★', '');
CREATE TABLE questiongroup_tags (
    Tag text NOT NULL,
    IdQuestiongroup integer NOT NULL
);
-- insert the common tags
INSERT INTO questiongroup_tags SELECT DISTINCT
    tag,
    Q1.idgroup
FROM
    question_tags AS Tag1
    JOIN questions2 AS Q1 ON Q1.id = Tag1.idquestion
WHERE
    Q1.idgroup IS NOT NULL
    AND (
        SELECT
            bool_and(Tag1.tag IN (
                    SELECT
                        tag FROM question_tags
                    WHERE
                        question_tags.idquestion = Q2.id))
        FROM
            questions2 AS Q2
        WHERE
            Q2.idgroup = Q1.idgroup);
-- now copy the exclusive tags to the subtitle
-- start by removing titles

UPDATE
    questions2
SET
    page = page - 'title';
UPDATE
    questions2
SET
    Subtitle = ''
WHERE
    NeedExercice IS NULL;
WITH tagstrings (
    idquestion,
    tags
) AS (
    SELECT
        question_tags.idquestion,
        string_agg(tag, ', ')
    FROM
        question_tags
        JOIN questions2 ON question_tags.idquestion = questions2.id
    WHERE
        tag NOT IN (
            SELECT
                questiongroup_tags.tag
            FROM
                questiongroup_tags
            WHERE
                questiongroup_tags.idquestiongroup = questions2.idgroup)
        GROUP BY
            question_tags.idquestion)
UPDATE
    questions2
SET
    Subtitle = tagstrings.tags
FROM
    tagstrings
WHERE
    tagstrings.idquestion = questions2.id;
-- remove the useless question_tags table
DROP TABLE question_tags;
-- temporary remove constraint on exercice_question
ALTER TABLE exercice_questions
    DROP CONSTRAINT exercice_questions_id_question_fkey;
-- switch to questions v2
DROP TABLE questions;
ALTER TABLE questions2 RENAME TO questions;
ALTER TABLE questions
    ADD PRIMARY KEY (id);
-- add the constraint back
ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;
-- temporary remove primary key, which is added back with the constraints
ALTER TABLE exercice_questions
    DROP CONSTRAINT exercice_questions_pkey CASCADE;
COMMIT;

--
-- exercices
--

BEGIN;
-- create an empty table
CREATE TABLE exercicegroup_tags (
    Tag text NOT NULL,
    IdExercicegroup integer NOT NULL
);
CREATE TABLE exercicegroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);
-- exercices : compute implicit groups with the title
INSERT INTO exercicegroups (title, public, idteacher)
SELECT DISTINCT
    title,
    public,
    idteacher
FROM
    exercices;
-- attach exercices to groups, by creating a new table, required to insert idgroup at the correct place
CREATE TABLE exercices2 (
    Id serial,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Parameters jsonb NOT NULL
);
INSERT INTO exercices2
SELECT
    id,
    (
        SELECT
            exercicegroups.id
        FROM
            exercicegroups
        WHERE
            exercicegroups.title = exercices.title),
    -- there is no tags for exercice yet, so the subtitle is empty
    '', Description, Parameters
FROM
    exercices;
-- we loose current progressions and tasks but it is OK here // TODO: maybe not, check that
ALTER TABLE exercice_questions
    DROP CONSTRAINT exercice_questions_id_exercice_fkey;
ALTER TABLE tasks
    DROP CONSTRAINT tasks_idexercice_fkey;
DROP TABLE exercices CASCADE;
ALTER TABLE exercices2 RENAME TO exercices;
ALTER TABLE exercices
    ADD PRIMARY KEY (id);
-- check by puting the constraint back
ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices ON DELETE CASCADE;
ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;
COMMIT;

--
-- Trivial : add explicit field for difficulty query
--
--
-- we assume difficulties are regular, meaning that
-- if one difficulty is specified in a categorie, it is for all

CREATE OR REPLACE FUNCTION __migration_has_difficulty (difficulty text, questions jsonb)
    RETURNS boolean
    AS $$
DECLARE
BEGIN
    RETURN (
        SELECT
            bool_and((
                SELECT
                    bool_or(intersection ? difficulty)
                    FROM jsonb_array_elements(
                        CASE WHEN union_ = 'null'::jsonb THEN
                            '[]'
                        ELSE
                            union_
                        END) AS intersection))
        FROM
            jsonb_array_elements(jsonb_strip_nulls (questions)) AS union_);
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_remove_diff (questions jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN (
        SELECT
            jsonb_agg((
                SELECT
                    jsonb_agg(DISTINCT __migration_sort_json_array (intersection - '★' - '★★' - '★★★'))
                    FROM jsonb_array_elements(
                        CASE WHEN union_ = 'null'::jsonb THEN
                            '[]'
                        ELSE
                            union_
                        END) AS intersection))
        FROM
            jsonb_array_elements(questions) AS union_);
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

CREATE OR REPLACE FUNCTION __migration_sort_json_array (array_ jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN (
        SELECT
            jsonb_agg(value ORDER BY value)
        FROM
            jsonb_array_elements(array_));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;

-- Step 1 : replace the questions array by a struct
BEGIN;
ALTER TABLE trivials
    DROP CONSTRAINT Questions_gomacro;
UPDATE
    trivials
SET
    Questions = jsonb_build_object('Tags', Questions, 'Difficulties', '[]'::jsonb);
-- Step 2 : fill the Difficulties array where needed
UPDATE
    trivials
SET
    Questions = jsonb_set(Questions, '{Difficulties}', Questions -> 'Difficulties' || '["★"]'::jsonb)
WHERE
    __migration_has_difficulty ('★', Questions -> 'Tags');
UPDATE
    trivials
SET
    Questions = jsonb_set(Questions, '{Difficulties}', Questions -> 'Difficulties' || '["★★"]'::jsonb)
WHERE
    __migration_has_difficulty ('★★', Questions -> 'Tags');
UPDATE
    trivials
SET
    Questions = jsonb_set(Questions, '{Difficulties}', Questions -> 'Difficulties' || '["★★★"]'::jsonb)
WHERE
    __migration_has_difficulty ('★★★', Questions -> 'Tags');
--
-- Step 3 : remove difficulty tags, removing induced duplicates

UPDATE
    trivials
SET
    Questions = jsonb_set(Questions, '{Tags}', __migration_remove_diff (Questions -> 'Tags'));
COMMIT;

BEGIN;
-- add monoquestion field
CREATE TABLE monoquestions (
    Id serial PRIMARY KEY,
    IdQuestion integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL
);
ALTER TABLE tasks
    ADD COLUMN IdMonoquestion integer;
-- remove idexercice from tasks
ALTER TABLE progressions
    DROP COLUMN idexercice CASCADE;
ALTER TABLE progression_questions
    DROP COLUMN idexercice CASCADE;
-- temporary remove primary key, which is added back with the constraints
ALTER TABLE sheet_tasks
    DROP CONSTRAINT sheet_tasks_pkey CASCADE;
COMMIT;

ALTER SEQUENCE questions2_id_seq
    RENAME TO questions_id_seq;

ALTER SEQUENCE exercices2_id_seq
    RENAME TO exercices_id_seq;

SELECT
    setval('questions_id_seq', (
            SELECT
                MAX(id)
            FROM questions));

SELECT
    setval('exercices_id_seq', (
            SELECT
                MAX(id)
            FROM exercices));

-- TODO: check standalone questions included in exercice
