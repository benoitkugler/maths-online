-- v1.2.3
-- 1) For clarity, slit up block list and parameters
-- 2) make sure only shared parameters are used in exercices
-- 3) move Description to Parameters

BEGIN;
-- create the new table
CREATE TABLE questions_2 (
    Id serial PRIMARY KEY,
    Ennonce jsonb NOT NULL,
    Parameters jsonb NOT NULL,
    Subtitle text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer
);
-- copy the old content in the new
INSERT INTO questions_2
SELECT
    Id,
    coalesce(page -> 'ennonce', '[]'),
    jsonb_set(coalesce(page -> 'parameters', '{}'), '{Description}', to_jsonb (Description)),
    Subtitle,
    Difficulty,
    NeedExercice,
    IdGroup
FROM
    questions;
-- same for exerices
CREATE TABLE exercices_2 (
    Id serial PRIMARY KEY,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    Parameters jsonb NOT NULL
);
-- copy the old content in the new
INSERT INTO exercices_2
SELECT
    Id,
    IdGroup,
    Subtitle,
    Difficulty,
    jsonb_set(coalesce(Parameters, '{}'), '{Description}', to_jsonb (Description))
FROM
    exercices;
--
--

SELECT
    id,
    NeedExercice,
    Parameters
FROM
    questions_2
WHERE
    NeedExercice IS NOT NULL;
ROLLBACK;

