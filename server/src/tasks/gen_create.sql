-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE progressions (
    Id serial PRIMARY KEY,
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL,
    IdExercice integer NOT NULL
);

CREATE TABLE progression_questions (
    IdProgression integer NOT NULL,
    IdExercice integer NOT NULL,
    Index integer NOT NULL,
    History boolean[]
);

CREATE TABLE tasks (
    Id serial PRIMARY KEY,
    IdExercice integer NOT NULL
);

-- constraints
ALTER TABLE tasks
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask);

ALTER TABLE progressions
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask, IdExercice) REFERENCES tasks (Id, IdExercice);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdExercice, INDEX) REFERENCES exercice_questions ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdProgression, IdExercice) REFERENCES progressions (Id, IdExercice) ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdProgression) REFERENCES progressions ON DELETE CASCADE;

ALTER TABLE progression_questions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices ON DELETE CASCADE;

