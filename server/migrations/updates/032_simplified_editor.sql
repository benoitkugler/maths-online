-- v1.3
ALTER TABLE teachers
    ADD COLUMN HasSimplifiedEditor boolean;

UPDATE
    teachers
SET
    HasSimplifiedEditor = FALSE;

ALTER TABLE teachers
    ALTER COLUMN HasSimplifiedEditor SET NOT NULL;

