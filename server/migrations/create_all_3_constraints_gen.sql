ALTER TABLE teachers
    ADD UNIQUE (Mail);

ALTER TABLE classrooms
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE students
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms ON DELETE CASCADE;

ALTER TABLE questions
    ADD CHECK (NeedExercice IS NOT NULL
        OR IdGroup IS NOT NULL);

ALTER TABLE questions
    ADD UNIQUE (Id, NeedExercice);

ALTER TABLE questions
    ADD FOREIGN KEY (NeedExercice) REFERENCES exercices;

ALTER TABLE questions
    ADD FOREIGN KEY (IdGroup) REFERENCES questiongroups ON DELETE CASCADE;

ALTER TABLE questiongroups
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE questiongroup_tags
    ADD UNIQUE (IdQuestiongroup, Tag);

ALTER TABLE questiongroup_tags
    ADD FOREIGN KEY (IdQuestiongroup) REFERENCES questiongroups ON DELETE CASCADE;

ALTER TABLE exercicegroups
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE exercicegroup_tags
    ADD UNIQUE (IdExercicegroup, Tag);

ALTER TABLE exercicegroup_tags
    ADD FOREIGN KEY (IdExercicegroup) REFERENCES exercicegroups ON DELETE CASCADE;

ALTER TABLE exercices
    ADD FOREIGN KEY (IdGroup) REFERENCES exercicegroups;

ALTER TABLE exercice_questions
    ADD PRIMARY KEY (IdExercice, INDEX);

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdExercice, IdQuestion) REFERENCES Questions (NeedExercice, Id);

ALTER TABLE exercice_questions
    ADD UNIQUE (IdQuestion);

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdExercice) REFERENCES exercices ON DELETE CASCADE;

ALTER TABLE exercice_questions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questions;

ALTER TABLE exercices
    ADD CONSTRAINT Parameters_gomacro CHECK (gomacro_validate_json_ques_Parameters (Parameters));

ALTER TABLE questions
    ADD CONSTRAINT Page_gomacro CHECK (gomacro_validate_json_ques_QuestionPage (Page));

ALTER TABLE trivials
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;

ALTER TABLE trivials
    ADD CONSTRAINT Questions_gomacro CHECK (gomacro_validate_json_triv_CategoriesQuestions (Questions));

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

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind);

ALTER TABLE review_questions
    ADD CHECK (Kind = 0
    /* Kind.KQuestion */);

ALTER TABLE review_questions
    ADD UNIQUE (IdQuestion);

ALTER TABLE review_questions
    ADD UNIQUE (IdReview);

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_questions
    ADD FOREIGN KEY (IdQuestion) REFERENCES questiongroups;

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind);

ALTER TABLE review_exercices
    ADD CHECK (Kind = 1
    /* Kind.KExercice */);

ALTER TABLE review_exercices
    ADD UNIQUE (IdExercice);

ALTER TABLE review_exercices
    ADD UNIQUE (IdReview);

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_exercices
    ADD FOREIGN KEY (IdExercice) REFERENCES exercicegroups;

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdReview, Kind) REFERENCES reviews (ID, Kind);

ALTER TABLE review_trivials
    ADD CHECK (Kind = 2
    /* Kind.KTrivial */);

ALTER TABLE review_trivials
    ADD UNIQUE (IdTrivial);

ALTER TABLE review_trivials
    ADD UNIQUE (IdReview);

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_trivials
    ADD FOREIGN KEY (IdTrivial) REFERENCES trivials;

ALTER TABLE review_participations
    ADD UNIQUE (IdReview, IdTeacher);

ALTER TABLE review_participations
    ADD FOREIGN KEY (IdReview) REFERENCES reviews ON DELETE CASCADE;

ALTER TABLE review_participations
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers ON DELETE CASCADE;

ALTER TABLE review_participations
    ADD CONSTRAINT Comments_gomacro CHECK (gomacro_validate_json_array_revi_Comment (Comments));

