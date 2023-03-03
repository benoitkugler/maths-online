-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE classrooms (
    Id serial PRIMARY KEY,
    IdTeacher integer NOT NULL,
    Name text NOT NULL
);

CREATE TABLE students (
    Id serial PRIMARY KEY,
    Name text NOT NULL,
    Surname text NOT NULL,
    Birthday date NOT NULL,
    TrivialSuccess integer NOT NULL,
    IsClientAttached boolean NOT NULL,
    IdClassroom integer NOT NULL
);

CREATE TABLE teachers (
    Id serial PRIMARY KEY,
    Mail text NOT NULL,
    PasswordCrypted bytea NOT NULL,
    IsAdmin boolean NOT NULL,
    HasSimplifiedEditor boolean NOT NULL
);

-- constraints
ALTER TABLE teachers
    ADD UNIQUE (Mail);

ALTER TABLE classrooms
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE students
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

