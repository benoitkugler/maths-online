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
    IsAdmin boolean NOT NULL
);

CREATE TABLE exercices (
    Id serial PRIMARY KEY,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Parameters jsonb NOT NULL
);

CREATE TABLE exercice_questions (
    IdExercice integer NOT NULL,
    IdQuestion integer NOT NULL,
    Bareme integer NOT NULL,
    Index integer NOT NULL
);

CREATE TABLE exercicegroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE exercicegroup_tags (
    Tag text NOT NULL,
    IdExercicegroup integer NOT NULL
);

CREATE TABLE questions (
    Id serial PRIMARY KEY,
    Page jsonb NOT NULL,
    Subtitle text NOT NULL,
    Description text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer
);

CREATE TABLE questiongroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE questiongroup_tags (
    Tag text NOT NULL,
    IdQuestiongroup integer NOT NULL
);

CREATE TABLE trivials (
    Id serial PRIMARY KEY,
    Questions jsonb NOT NULL,
    QuestionTimeout integer NOT NULL,
    ShowDecrassage boolean NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL,
    Name text NOT NULL
);

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
