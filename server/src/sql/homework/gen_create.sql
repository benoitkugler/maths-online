-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    IdTeacher integer NOT NULL,
    Level text NOT NULL,
    Anonymous integer,
    Public boolean NOT NULL,
    Matiere text CHECK (Matiere IN ('ALLEMAND', 'ANGLAIS', 'AUTRE', 'ESPAGNOL', 'FRANCAIS', 'HISTOIRE-GEO', 'ITALIEN', 'MATHS', 'PHYSIQUE', 'SES', 'SVT')) NOT NULL
);

CREATE TABLE sheet_tasks (
    IdSheet integer NOT NULL,
    Index integer NOT NULL,
    IdTask integer NOT NULL
);

CREATE TABLE travails (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    IdSheet integer NOT NULL,
    Noted boolean NOT NULL,
    Deadline timestamp(0) with time zone NOT NULL
);

-- constraints
ALTER TABLE travails
    ADD UNIQUE (Id, IdSheet);

ALTER TABLE travails
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE travails
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (Id, Anonymous) REFERENCES travails (IdSheet, Id) ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE sheets
    ADD FOREIGN KEY (Anonymous) REFERENCES travails ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD PRIMARY KEY (IdSheet, INDEX);

ALTER TABLE sheet_tasks
    ADD UNIQUE (IdTask);

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdTask) REFERENCES tasks;

