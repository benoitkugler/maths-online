-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE monoquestions (
    Id serial PRIMARY KEY,
    IdQuestion integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL
);

CREATE TABLE progressions (
    IdStudent integer NOT NULL,
    IdTask integer NOT NULL,
    Index integer NOT NULL,
    History boolean[]
);

CREATE TABLE random_monoquestions (
    Id serial PRIMARY KEY,
    IdQuestiongroup integer NOT NULL,
    NbRepeat integer NOT NULL,
    Bareme integer NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL
);

CREATE TABLE random_monoquestion_variants (
    IdStudent integer NOT NULL,
    IdRandomMonoquestion integer NOT NULL,
    Index integer NOT NULL,
    IdQuestion integer NOT NULL
);

CREATE TABLE tasks (
    Id serial PRIMARY KEY,
    IdExercice integer,
    IdMonoquestion integer,
    IdRandomMonoquestion integer
);

-- constraints
ALTER TABLE monoquestions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE random_monoquestions
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups;

ALTER TABLE tasks
    ADD UNIQUE (Id, IdExercice);

ALTER TABLE tasks
    ADD CHECK ((IdExercice IS NOT NULL)::int + (IdMonoquestion IS NOT NULL)::int + (IdRandomMonoquestion IS NOT NULL)::int = 1);

ALTER TABLE tasks
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdMonoquestion) REFERENCES monoquestions;

ALTER TABLE tasks
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions;

ALTER TABLE random_monoquestion_variants
    ADD UNIQUE (IdStudent, IdRandomMonoquestion, INDEX);

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdStudent) REFERENCES students;

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions;

ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE progressions
    ADD UNIQUE (IdStudent, IdTask, INDEX);

ALTER TABLE progressions
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;

ALTER TABLE progressions
    ADD FOREIGN KEY (IdTask) REFERENCES tasks ON DELETE CASCADE;

