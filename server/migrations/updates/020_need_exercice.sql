--v0.6.3
-- use the id of the exercice instead of a vague boolean type

BEGIN;
CREATE TABLE _tmp (
    id_question integer,
    id_exercice integer
);
INSERT INTO _tmp
SELECT
    id,
    (
        CASE WHEN need_exercice = TRUE THEN
        (
            SELECT
                exercice_questions.id_exercice
            FROM
                exercice_questions
            WHERE
                exercice_questions.id_question = questions.id)
        ELSE
            NULL
        END)
FROM
    questions;
ALTER TABLE questions
    ALTER COLUMN need_exercice DROP NOT NULL;
ALTER TABLE questions
    ALTER COLUMN need_exercice TYPE integer
    USING NULL;
UPDATE
    questions
SET
    need_exercice = _tmp.id_exercice
FROM
    _tmp
WHERE
    _tmp.id_question = questions.id;
DROP TABLE _tmp;
COMMIT;

