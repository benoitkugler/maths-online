-- v1.4
BEGIN;
--
ALTER TABLE questions
    ADD COLUMN Correction jsonb;
UPDATE
    questions
SET
    Correction = '[]';
-- empty list
ALTER TABLE questions
    ALTER COLUMN Correction SET NOT NULL;
ALTER TABLE questions
    ADD CONSTRAINT Correction_gomacro CHECK (gomacro_validate_json_array_ques_Block (Correction));
--
--

COMMIT;

