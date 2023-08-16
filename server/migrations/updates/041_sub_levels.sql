-- v1.6.1
-- add Spé to 1ère
-- check no exercices are concerned

SELECT
    count(*)
FROM
    exercicegroup_tags
WHERE
    tag = '1ERE';

-- classify all admin questions in 1ERE as SPE
BEGIN;
-- update constraint
ALTER TABLE exercicegroup_tags
    DROP CONSTRAINT section_enum;
ALTER TABLE exercicegroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (2, 1, 4, 3));
--
ALTER TABLE questiongroup_tags
    DROP CONSTRAINT section_enum;
ALTER TABLE questiongroup_tags
    ADD CONSTRAINT section_enum CHECK (Section IN (2, 1, 4, 3));
--
INSERT INTO questiongroup_tags (tag, section, idquestiongroup)
SELECT
    'SPE',
    4
    /** SubLevel */
,
    questiongroups.id
FROM
    questiongroup_tags
    JOIN questiongroups ON questiongroups.id = questiongroup_tags.idquestiongroup
WHERE
    questiongroups.idteacher = 1
    AND tag = '1ERE';
--
COMMIT;

