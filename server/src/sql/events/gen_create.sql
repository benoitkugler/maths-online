-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE events (
    IdStudent integer NOT NULL,
    Event smallint CHECK (Event IN (0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)) NOT NULL,
    Date date NOT NULL
);

-- constraints
ALTER TABLE events
    ADD FOREIGN KEY (IdStudent) REFERENCES students;

