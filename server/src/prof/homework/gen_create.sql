-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    Title text NOT NULL,
    Notation integer CHECK (Notation IN (0, 1)) NOT NULL,
    Activated boolean NOT NULL,
    Deadline timestamp(0) with time zone NOT NULL
);

CREATE TABLE sheet_exercices (
    IdSheet integer NOT NULL,
    IdExercice integer NOT NULL,
    Index integer NOT NULL
);

CREATE TABLE student_progressions (
    IdStudent integer NOT NULL,
    IdSheet integer NOT NULL,
    Index integer NOT NULL,
    IdExercice integer NOT NULL,
    IdProgression integer NOT NULL
);

-- constraints
ALTER TABLE sheets
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE sheet_exercices
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheet_exercices
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdProgression) REFERENCES progressions;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdProgression, IdExercice) REFERENCES Progression (Id, IdExercice) ON DELETE CASCADE;

ALTER TABLE student_progressions
    ADD FOREIGN KEY (IdSheet, IdExercice, INDEX) REFERENCES sheet_exercices ON DELETE CASCADE;

ALTER TABLE student_progressions
    ADD UNIQUE (IdStudent, IdSheet, INDEX);

