-- v1.2.0
-- refactor tags handling by introducing sections
-- 1 : 'Niveau'
-- 2 : 'Chapitre'
-- 3 : 'TrivMaths'

BEGIN;
--
--
--

ALTER TABLE questiongroup_tags
    ADD COLUMN Section integer;
ALTER TABLE exercicegroup_tags
    ADD COLUMN Section integer;
UPDATE
    questiongroup_tags
SET
    Section = 3;
UPDATE
    exercicegroup_tags
SET
    Section = 3;
--
-- Attribute level section

UPDATE
    questiongroup_tags
SET
    Section = 1
WHERE
    tag = ANY (ARRAY['2NDE', '1ERE', 'TERM', 'CPGE']);
UPDATE
    exercicegroup_tags
SET
    Section = 1
WHERE
    tag = ANY (ARRAY['2NDE', '1ERE', 'TERM', 'CPGE']);
-- Attribute chapters section
UPDATE
    questiongroup_tags
SET
    Section = 2
WHERE
    tag = ANY (ARRAY['VECTEURS', 'EXEMPLE', 'DERIVATION', 'SUITES', 'VARIATION D''UNE FONCTION', 'REPERAGE DU PLAN', 'DEGRE 2', 'DROITES', 'PRODUIT SCALAIRE', 'REELS', 'ENTIERS', 'POURCENTAGES', 'CEINTURE BLANCHE', 'CEINTURE JAUNE', 'CEINTURE ORANGE', 'CEINTURE BLEUE', 'CEINTURE VERTE', 'GENERALITES SUR LES FONCTIONS', 'STATISTIQUES']);
UPDATE
    exercicegroup_tags
SET
    Section = 2
WHERE
    tag = ANY (ARRAY['EXEMPLE']);
--
-- add constraints

ALTER TABLE questiongroup_tags
    ALTER COLUMN Section SET NOT NULL;
ALTER TABLE exercicegroup_tags
    ALTER COLUMN Section SET NOT NULL;
ALTER TABLE questiongroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (1, 2, 3));
ALTER TABLE exercicegroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (1, 2, 3));
-- check upper case
ALTER TABLE questiongroup_tags
    ADD CONSTRAINT tag_upper CHECK (tag = upper(tag));
ALTER TABLE exercicegroup_tags
    ADD CONSTRAINT tag_upper CHECK (tag = upper(tag));
-- enfore unique level and chapter
CREATE UNIQUE INDEX questiongroup_tags_level ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 1;
CREATE UNIQUE INDEX questiongroup_tags_chapter ON questiongroup_tags (IdQuestiongroup)
WHERE
    Section = 2;
CREATE UNIQUE INDEX exercicegroup_tags_level ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 1;
CREATE UNIQUE INDEX exercicegroup_tags_chapter ON exercicegroup_tags (IdExercicegroup)
WHERE
    Section = 2;
--
-- update the trivial query tags, by adding the matching Section
--

CREATE OR REPLACE FUNCTION __migration_build_TagSection (inTag text)
    RETURNS jsonb
    AS $$
DECLARE
BEGIN
    -- search for the matching section
    -- we assume there is only one section for a given tag
    RETURN jsonb_build_object('Tag', upper(inTag), 'Section', (
            SELECT
                section FROM questiongroup_tags
            WHERE
                questiongroup_tags.tag = upper(inTag)
        LIMIT 1));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
--
--

UPDATE
    trivials
SET
    questions = jsonb_set(questions, '{Tags}', (
            SELECT
                jsonb_agg((
                    SELECT
                        jsonb_agg((
                            SELECT
                                jsonb_agg(__migration_build_TagSection (tag #>> '{}'))
                        FROM jsonb_array_elements(intersection) AS tag))
            FROM jsonb_array_elements(categorie) AS intersection))
    FROM jsonb_array_elements(questions -> 'Tags') AS categorie));
DROP FUNCTION __migration_build_tagsection (text);
COMMIT;

