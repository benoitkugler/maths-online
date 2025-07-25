BEGIN;
CREATE OR REPLACE FUNCTION __migration_int_array (Advance jsonb)
    RETURNS smallint[]
    AS $$
DECLARE
BEGIN
    RETURN (
        SELECT
            array_agg(value::smallint)
        FROM
            jsonb_array_elements(Advance));
END;
$$
LANGUAGE 'plpgsql'
IMMUTABLE;
-- enum types
ALTER TABLE exercicegroup_tags
    ALTER COLUMN Section TYPE smallint;
ALTER TABLE questiongroup_tags
    ALTER COLUMN Section TYPE smallint;
ALTER TABLE reviews
    ALTER COLUMN Kind TYPE smallint;
ALTER TABLE review_exercices
    ALTER COLUMN Kind TYPE smallint;
ALTER TABLE review_questions
    ALTER COLUMN Kind TYPE smallint;
ALTER TABLE review_sheets
    ALTER COLUMN Kind TYPE smallint;
ALTER TABLE review_trivials
    ALTER COLUMN Kind TYPE smallint;
ALTER TABLE review_participations
    ALTER COLUMN Approval TYPE smallint;
ALTER TABLE beltevolutions
    ALTER COLUMN Level TYPE smallint;
ALTER TABLE beltquestions
    ALTER COLUMN DOMAIN TYPE smallint;
ALTER TABLE beltquestions
    ALTER COLUMN Rank TYPE smallint;
-- date
ALTER TABLE students
    ALTER COLUMN Birthday TYPE date;
-- int array
ALTER TABLE beltevolutions
    DROP CONSTRAINT advance_gomacro;
ALTER TABLE beltevolutions
    ALTER COLUMN Advance TYPE smallint[]
    USING (__migration_int_array (Advance));
ALTER TABLE beltevolutions
    ALTER COLUMN Advance SET NOT NULL;
ALTER TABLE beltevolutions
    ADD CHECK (array_length(Advance, 1) = 12);
DROP FUNCTION __migration_int_array;
COMMIT;

