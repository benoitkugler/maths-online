--  BUG FIX
ALTER TABLE events
    DROP CONSTRAINT events_idstudent_fkey;

ALTER TABLE events
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

