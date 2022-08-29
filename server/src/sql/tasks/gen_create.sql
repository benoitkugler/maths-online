-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE monoquestions (
    Id serial PRIMARY KEY,
    IdQuestion integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL
);

CREATE TABLE progressions (
    Id serial PRIMARY KEY,
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL
);

CREATE TABLE progression_questions (
    IdProgression integer NOT NULL,
    Index integer NOT NULL,
    History boolean[]
);

CREATE TABLE tasks (
    Id serial PRIMARY KEY,
    IdExercice integer,
    IdMonoquestion integer
);

-- constraints
ALTER TABLE monoquestions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE tasks
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE tasks
    ADD CHECK (IdExercice IS NOT NULL
        OR IdMonoquestion IS NOT NULL);

ALTER TABLE tasks
    ADD CHECK (IdExercice IS NULL
        OR IdMonoquestion IS NULL);

ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdMonoquestion) REFERENCES monoquestions;

ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdProgression) REFERENCES progressions ON DELETE CASCADE;

