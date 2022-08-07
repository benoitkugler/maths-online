--v0.6.3
-- update the column names to the new gomacro code generator

BEGIN;
ALTER TABLE classrooms RENAME COLUMN id_teacher TO IdTeacher;
ALTER TABLE students RENAME COLUMN id_classroom TO IdClassroom;
ALTER TABLE teachers RENAME COLUMN password_crypted TO PasswordCrypted;
ALTER TABLE teachers RENAME COLUMN is_admin TO IsAdmin;
ALTER TABLE exercices RENAME COLUMN id_teacher TO IdTeacher;
ALTER TABLE exercice_questions RENAME COLUMN id_exercice TO IdExercice;
ALTER TABLE exercice_questions RENAME COLUMN id_question TO IdQuestion;
ALTER TABLE progressions RENAME COLUMN id_exercice TO IdExercice;
ALTER TABLE progression_questions RENAME COLUMN id_progression TO IdProgression;
ALTER TABLE progression_questions RENAME COLUMN id_exercice TO IdExercice;
ALTER TABLE questions RENAME COLUMN id_teacher TO IdTeacher;
ALTER TABLE questions RENAME COLUMN need_exercice TO NeedExercice;
ALTER TABLE question_tags RENAME COLUMN id_question TO IdQuestion;
ALTER TABLE trivial_configs RENAME TO trivials;
ALTER TABLE trivials RENAME COLUMN id_teacher TO IdTeacher;
COMMIT;

