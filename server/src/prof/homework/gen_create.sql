-- Code genererated by gomacro/generator/sql. DO NOT EDIT.
CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    IdClassroom integer NOT NULL,
    Title text NOT NULL,
    Notation integer CHECK (Notation IN (0, 1)) NOT NULL,
    Activated boolean NOT NULL,
    Deadline timestamp(0) with time zone NOT NULL
);

CREATE TABLE sheet_tasks (
    IdSheet integer NOT NULL,
    Index integer NOT NULL,
    IdTask integer NOT NULL
);

-- constraints
ALTER TABLE sheets
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD PRIMARY KEY (IdSheet, INDEX);

ALTER TABLE sheet_tasks
    ADD UNIQUE (IdTask);

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdSheet) REFERENCES sheets ON DELETE CASCADE;

ALTER TABLE sheet_tasks
    ADD FOREIGN KEY (IdTask) REFERENCES tasks;

