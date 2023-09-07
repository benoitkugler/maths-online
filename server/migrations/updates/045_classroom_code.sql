-- persists classroom code on DB
BEGIN;
--
CREATE TABLE classroom_codes (
    IdClassroom integer NOT NULL,
    Code text NOT NULL,
    ExpiresAt timestamp(0
) with time zone NOT NULL
);
--
ALTER TABLE classroom_codes
    ADD UNIQUE (Code);
ALTER TABLE classroom_codes
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;
--
COMMIT;

