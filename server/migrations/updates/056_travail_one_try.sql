BEGIN;
ALTER TABLE travails
    ADD COLUMN QuestionRepeat smallint CHECK (QuestionRepeat IN (0, 1));
UPDATE
    travails
SET
    QuestionRepeat = 0;
ALTER TABLE travails
    ALTER COLUMN QuestionRepeat SET NOT NULL;
ALTER TABLE travails
    ADD COLUMN QuestionTimeLimit integer;
UPDATE
    travails
SET
    QuestionTimeLimit = 0;
ALTER TABLE travails
    ALTER COLUMN QuestionTimeLimit SET NOT NULL;
COMMIT;

