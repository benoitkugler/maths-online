-- v0.4.4
ALTER TABLE questions
    ADD COLUMN description varchar;

UPDATE
    questions
SET
    description = '';

ALTER TABLE questions
    ALTER COLUMN description SET NOT NULL;

