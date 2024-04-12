-- v1.6.3
-- add ShowAfter time field

BEGIN;
ALTER TABLE travails
    ADD COLUMN ShowAfter timestamp(0) WITH time zone;
UPDATE
    travails
SET
    ShowAfter = NOW();
ALTER TABLE travails
    ALTER COLUMN ShowAfter SET NOT NULL;
COMMIT;

