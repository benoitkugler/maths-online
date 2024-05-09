BEGIN;
ALTER TABLE beltquestions
    ADD COLUMN Title text;
UPDATE
    beltquestions
SET
    Title = '';
ALTER TABLE beltquestions
    ALTER COLUMN Title SET NOT NULL;
COMMIT;

