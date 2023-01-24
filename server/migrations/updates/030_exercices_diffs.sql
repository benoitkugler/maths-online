-- v1.2.3
-- add Difficulty field to exercices

BEGIN;
ALTER TABLE exercices
    ADD COLUMN Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', ''));
--
UPDATE
    exercices
SET
    Difficulty = '';
--
ALTER TABLE exercices
    ALTER COLUMN Difficulty SET NOT NULL;
COMMIT;

