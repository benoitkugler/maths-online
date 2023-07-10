BEGIN;
ALTER TABLE teachers
    ADD COLUMN Contact jsonb;
UPDATE
    teachers
SET
    Contact = '{"Name":"", "URL":""}';
ALTER TABLE teachers
    ALTER COLUMN Contact SET NOT NULL;
ALTER TABLE teachers
    ADD CONSTRAINT Contact_gomacro CHECK (gomacro_validate_json_teac_Contact (Contact));
COMMIT;

