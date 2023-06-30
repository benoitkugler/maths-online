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

CREATE TABLE exercices (
    Id serial PRIMARY KEY,
    IdGroup integer NOT NULL,
    Subtitle text NOT NULL,
    Parameters jsonb NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL
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
    IdExercicegroup integer NOT NULL,
    Section integer CHECK (Section IN (2, 1, 3)) NOT NULL
);

CREATE TABLE questions (
    Id serial PRIMARY KEY,
    Subtitle text NOT NULL,
    Difficulty text CHECK (Difficulty IN ('★', '★★', '★★★', '')) NOT NULL,
    NeedExercice integer,
    IdGroup integer,
    Enonce jsonb NOT NULL,
    Parameters jsonb NOT NULL,
    Correction jsonb NOT NULL
);

CREATE TABLE questiongroups (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    Public boolean NOT NULL,
    IdTeacher integer NOT NULL
);

CREATE TABLE questiongroup_tags (
    Tag text NOT NULL,
    IdQuestiongroup integer NOT NULL,
    Section integer CHECK (Section IN (2, 1, 3)) NOT NULL
);

CREATE TABLE selfaccess_trivials (
    IdClassroom integer NOT NULL,
    IdTrivial integer NOT NULL,
    IdTeacher integer NOT NULL
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
    Difficulty jsonb NOT NULL
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

CREATE TABLE sheets (
    Id serial PRIMARY KEY,
    Title text NOT NULL,
    IdTeacher integer NOT NULL,
    Level text NOT NULL
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

CREATE TABLE reviews (
    Id serial PRIMARY KEY,
    Kind integer CHECK (Kind IN (0, 1, 2)) NOT NULL
);

CREATE TABLE review_exercices (
    IdReview integer NOT NULL,
    IdExercice integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2)) NOT NULL
);

CREATE TABLE review_participations (
    IdReview integer NOT NULL,
    IdTeacher integer NOT NULL,
    Approval integer CHECK (Approval IN (0, 1, 2)) NOT NULL,
    Comments jsonb NOT NULL
);

CREATE TABLE review_questions (
    IdReview integer NOT NULL,
    IdQuestion integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2)) NOT NULL
);

CREATE TABLE review_trivials (
    IdReview integer NOT NULL,
    IdTrivial integer NOT NULL,
    Kind integer CHECK (Kind IN (0, 1, 2)) NOT NULL
);

