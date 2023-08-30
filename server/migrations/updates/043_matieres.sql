-- add support for several topics
BEGIN;
ALTER TABLE teachers
    ADD COLUMN FavoriteMatiere text;
--
UPDATE
    teachers
SET
    FavoriteMatiere = 'MATHS';
--
ALTER TABLE teachers
    ALTER COLUMN FavoriteMatiere SET NOT NULL;
--
ALTER TABLE teachers
    ADD CONSTRAINT FavoriteMatiere_check CHECK (FavoriteMatiere IN ('ALLEMAND', 'ANGLAIS', 'AUTRE', 'ESPAGNOL', 'FRANCAIS', 'HISTOIRE-GEO', 'ITALIEN', 'MATHS', 'PHYSIQUE', 'SES', 'SVT'));
--
-- update section enum

ALTER TABLE exercicegroup_tags
    DROP CONSTRAINT section_enum;
ALTER TABLE exercicegroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (2, 1, 4, 3, 5));
--
ALTER TABLE questiongroup_tags
    DROP CONSTRAINT section_enum;
ALTER TABLE questiongroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (2, 1, 4, 3, 5));
-- add unique constraints
CREATE UNIQUE INDEX QuestiongroupTag_matiere ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 5
    /* Section.Matiere */
;
CREATE UNIQUE INDEX ExercicegroupTag_matiere ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 5
    /* Section.Matiere */
;
-- classify existing questions
INSERT INTO questiongroup_tags (IdQuestiongroup, Section, Tag)
SELECT
    Id,
    5,
    'MATHS'
FROM
    questiongroups
WHERE
    IdTeacher != 13;
INSERT INTO questiongroup_tags (IdQuestiongroup, Section, Tag)
SELECT
    Id,
    5,
    'AUTRE'
FROM
    questiongroups
WHERE
    IdTeacher = 13;
INSERT INTO exercicegroup_tags (IdExercicegroup, Section, Tag)
SELECT
    Id,
    5,
    'MATHS'
FROM
    exercicegroups
WHERE
    IdTeacher != 13;
INSERT INTO exercicegroup_tags (IdExercicegroup, Section, Tag)
SELECT
    Id,
    5,
    'AUTRE'
FROM
    exercicegroups
WHERE
    IdTeacher = 13;
-- update existing trivmaths
--

CREATE OR REPLACE FUNCTION __migration_add_matiere (questions jsonb, new_tags jsonb)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    RETURN jsonb_set(questions, '{Tags}', (
            SELECT
                jsonb_agg(
                    CASE WHEN categorie = 'null' THEN
                        '[]'
                    ELSE
                        (
                            SELECT
                                jsonb_agg(intersection || new_tags)
                            FROM jsonb_array_elements(categorie) AS intersection)
                    END)
                FROM jsonb_array_elements(questions -> 'Tags') AS categorie));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
UPDATE
    trivials
SET
    Questions = __migration_add_matiere (Questions, '[{"Tag": "MATHS", "Section": 5}]'::jsonb)
WHERE
    IdTeacher != 13;
UPDATE
    trivials
SET
    Questions = __migration_add_matiere (Questions, '[{"Tag": "AUTRE", "Section": 5}]'::jsonb)
WHERE
    IdTeacher = 13;
--
DROP FUNCTION __migration_add_matiere (jsonb, jsonb);
--
-- update Sheets

ALTER TABLE sheets
    ADD COLUMN Matiere text;
UPDATE
    sheets
SET
    Matiere = 'MATHS';
UPDATE
    sheets
SET
    Matiere = 'AUTRE'
WHERE
    idteacher = 13;
--
ALTER TABLE sheets
    ALTER COLUMN Matiere SET NOT NULL;
--
ALTER TABLE sheets
    ADD CONSTRAINT Matiere_check CHECK (Matiere IN ('ALLEMAND', 'ANGLAIS', 'AUTRE', 'ESPAGNOL', 'FRANCAIS', 'HISTOIRE-GEO', 'ITALIEN', 'MATHS', 'PHYSIQUE', 'SES', 'SVT'));
--
COMMIT;

