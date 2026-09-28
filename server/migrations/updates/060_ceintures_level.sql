BEGIN;

ALTER TABLE
    beltquestions
ADD
    COLUMN LevelMin smallint;

UPDATE
    beltquestions
SET
    LevelMin = 0;

ALTER TABLE
    beltquestions
ALTER COLUMN
    LevelMin
SET
    NOT NULL;

ALTER TABLE
    beltquestions
ADD
    CONSTRAINT beltquestions_levelmin_check CHECK (LevelMin IN (0, 1, 2, 3));

COMMIT;