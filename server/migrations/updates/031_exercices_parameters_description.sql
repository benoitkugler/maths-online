-- v1.2.3
-- 1) For clarity, slit up block list and parameters
-- 2) move Description to Parameters
-- 3) update Parameters to the new list form

BEGIN;
-- Step 1 : Add the new column...
ALTER TABLE questions
    ADD COLUMN Enonce jsonb;
ALTER TABLE questions
    ADD COLUMN Parameters jsonb;
-- ...  copying the content
UPDATE
    questions
SET
    Enonce = coalesce(page -> 'enonce', '[]');
--
UPDATE
    questions
SET
    Parameters = jsonb_set(coalesce(page -> 'parameters', '{}'), '{Description}', to_jsonb (Description));
-- add the NOT NULL constraint
ALTER TABLE questions
    ALTER COLUMN Enonce SET NOT NULL;
ALTER TABLE questions
    ALTER COLUMN Parameters SET NOT NULL;
-- Step 2 : DROP the unusued Description and Page columns
ALTER TABLE questions
    DROP COLUMN Description;
ALTER TABLE questions
    DROP COLUMN page;
--
--
-- Exercices
-- Step 0 : temporarily drop the constraint

ALTER TABLE exercices
    DROP CONSTRAINT parameters_gomacro;
-- Step 1 : update the Parameters column
UPDATE
    exercices
SET
    Parameters = jsonb_set(coalesce(Parameters, '{}'), '{Description}', to_jsonb (Description));
-- Step 2 : drop the description field
ALTER TABLE exercices
    DROP COLUMN Description;
--
UPDATE
    questions
SET
    Parameters = (
        CASE WHEN Parameters ->> 'Description' != '' THEN
            jsonb_build_array(jsonb_build_object('Kind', 'Co', 'Data', Parameters -> 'Description'))
        ELSE
            '[]'
        END) || coalesce((
        SELECT
            jsonb_agg(jsonb_build_object('Kind', 'Rp', 'Data', rp))
        FROM jsonb_array_elements(
            CASE WHEN coalesce(Parameters -> 'Variables', '[]')::text = 'null' THEN
                '[]'
            ELSE
                coalesce(Parameters -> 'Variables', '[]')
            END) AS rp), '[]') || coalesce((
        SELECT
            jsonb_agg(jsonb_build_object('Kind', 'In', 'Data', intr))
        FROM jsonb_array_elements(
            CASE WHEN coalesce(Parameters -> 'Intrinsics', '[]')::text = 'null' THEN
                '[]'
            ELSE
                coalesce(Parameters -> 'Intrinsics', '[]')
            END) AS intr), '[]');
UPDATE
    exercices
SET
    Parameters = (
        CASE WHEN Parameters ->> 'Description' != '' THEN
            jsonb_build_array(jsonb_build_object('Kind', 'Co', 'Data', Parameters -> 'Description'))
        ELSE
            '[]'
        END) || coalesce((
        SELECT
            jsonb_agg(jsonb_build_object('Kind', 'Rp', 'Data', rp))
        FROM jsonb_array_elements(
            CASE WHEN coalesce(Parameters -> 'Variables', '[]')::text = 'null' THEN
                '[]'
            ELSE
                coalesce(Parameters -> 'Variables', '[]')
            END) AS rp), '[]') || coalesce((
        SELECT
            jsonb_agg(jsonb_build_object('Kind', 'In', 'Data', intr))
        FROM jsonb_array_elements(
            CASE WHEN coalesce(Parameters -> 'Intrinsics', '[]')::text = 'null' THEN
                '[]'
            ELSE
                coalesce(Parameters -> 'Intrinsics', '[]')
            END) AS intr), '[]');
COMMIT;

