-- v1.6.4
CREATE TABLE travail_exceptions (
    IdStudent integer NOT NULL,
    IdTravail integer NOT NULL,
    Deadline timestamp(0
) with time zone,
    IgnoreForMark boolean NOT NULL
);

ALTER TABLE travail_exceptions
    ADD UNIQUE (IdStudent, IdTravail);

ALTER TABLE travail_exceptions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE travail_exceptions
    ADD FOREIGN KEY (IdTravail) REFERENCES travails ON DELETE CASCADE;

