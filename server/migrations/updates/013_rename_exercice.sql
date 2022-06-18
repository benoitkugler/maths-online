-- v0.5.1
BEGIN;
--
-- rename the package exercice to questions, requires editor/gen_create.sql
--

ALTER TABLE questions
    DROP CONSTRAINT question_structgen_validate_json_exe_questionpage;
ALTER TABLE questions
    ADD CONSTRAINT page_structgen_validate_json_que_QuestionPage CHECK (structgen_validate_json_que_QuestionPage (page));
--
-- add NeedExercice field
--

ALTER TABLE questions
    ADD COLUMN need_exercice boolean;
UPDATE
    questions
SET
    need_exercice = FALSE;
ALTER TABLE questions
    ALTER COLUMN need_exercice SET NOT NULL;
COMMIT;

