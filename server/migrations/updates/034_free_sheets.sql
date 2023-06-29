-- v1.5
-- split Sheet into Travail and Sheet

BEGIN;
CREATE TABLE travails (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    IdSheet integer NOT NULL,
    Noted boolean NOT NULL,
    Deadline timestamp(0
) with time zone NOT NULL
);
-- create [Travail] from [Sheet]
INSERT INTO travails (idclassroom, idsheet, noted, deadline)
SELECT
    idclassroom,
    id,
    notation = 1,
    CASE WHEN notation = 1 THEN
        deadline
    ELSE
        '0001-01-01 00:00:00Z'
    END
FROM
    sheets
WHERE
    activated = TRUE;
;
-- update sheets
-- map classrooms to their teacher

ALTER TABLE sheets
    ADD COLUMN IdTeacher integer;
UPDATE
    sheets
SET
    IdTeacher = classrooms.IdTeacher
FROM
    classrooms
WHERE
    sheets.idclassroom = classrooms.id;
-- map classrooms name to levels
CREATE OR REPLACE FUNCTION __migration_classroom_to_level (name text)
    RETURNS text
    AS $$
DECLARE
BEGIN
    RETURN CASE WHEN name LIKE '2G%' THEN
        '2NDE'
    WHEN name LIKE 'ECG1%' THEN
        'CPGE'
    ELSE
        ''
    END;
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
ALTER TABLE sheets
    ADD COLUMN Level text;
UPDATE
    sheets
SET
    Level = __migration_classroom_to_level (classrooms.Name)
FROM
    classrooms
WHERE
    sheets.idclassroom = classrooms.id;
ALTER TABLE sheets
    DROP COLUMN notation;
ALTER TABLE sheets
    DROP COLUMN idclassroom;
ALTER TABLE sheets
    DROP COLUMN deadline;
ALTER TABLE sheets
    DROP COLUMN activated;
SELECT
    *
FROM
    sheets;
SELECT
    *
FROM
    travails;
---
ALTER TABLE sheets
    ALTER COLUMN IdTeacher SET NOT NULL;
ALTER TABLE sheets
    ALTER COLUMN Level SET NOT NULL;
ALTER TABLE travails
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;
ALTER TABLE travails
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;
ALTER TABLE sheets
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;
COMMIT;

