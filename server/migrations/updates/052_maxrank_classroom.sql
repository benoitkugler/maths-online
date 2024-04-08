BEGIN;
ALTER TABLE classrooms
    ADD COLUMN MaxRankThreshold integer;
UPDATE
    classrooms
SET
    MaxRankThreshold = 40000;
ALTER TABLE classrooms
    ALTER COLUMN MaxRankThreshold SET NOT NULL;
COMMIT;

