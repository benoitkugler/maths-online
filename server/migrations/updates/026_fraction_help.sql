--v0.9.0
-- add ShowFractionHelp to ExpressionFieldBlock

BEGIN;
UPDATE
    questions
SET
    page = jsonb_set(page, '{enonce}', coalesce((
            SELECT
                jsonb_agg(
                    CASE WHEN value ->> 'Kind' = 'ExpressionFieldBlock' THEN
                        jsonb_set(value, '{Data, ShowFractionHelp}', 'false')
                ELSE
                    value
                    END)
                FROM jsonb_array_elements(page -> 'enonce')), '[]'));
COMMIT;

-- restore the constraint
ALTER TABLE questions
    ADD CONSTRAINT Page_gomacro CHECK (gomacro_validate_json_ques_QuestionPage (Page));

