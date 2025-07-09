-- delete progression on student suppression
BEGIN;
ALTER TABLE beltevolutions
    DROP CONSTRAINT beltevolutions_idstudent_fkey;
ALTER TABLE beltevolutions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;
COMMIT;

