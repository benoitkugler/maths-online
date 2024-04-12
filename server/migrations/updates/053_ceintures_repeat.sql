BEGIN;
ALTER TABLE beltquestions
    ADD COLUMN Repeat integer;
UPDATE
    beltquestions
SET
    Repeat = 1;
ALTER TABLE beltquestions
    ALTER COLUMN Repeat SET NOT NULL;
COMMIT;

