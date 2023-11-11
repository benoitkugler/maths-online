BEGIN;
ALTER TABLE students
    ADD COLUMN clients jsonb;
UPDATE
    students
SET
    clients = CASE WHEN IsClientAttached = TRUE THEN
        '[{"Device":"","Time":"2023-11-10T20:00:00.0+01:00"}]'::jsonb
    ELSE
        '[]'
    END;
ALTER TABLE students
    ALTER COLUMN clients SET NOT NULL;
ALTER TABLE students
    DROP COLUMN TrivialSuccess;
ALTER TABLE students
    DROP COLUMN IsClientAttached;
COMMIT;

