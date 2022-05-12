-- v0.2.5
-- renaming a block type changes the order and the auto-generated kind int value

UPDATE
    questions
SET
    enonce = (
        SELECT
            jsonb_agg(
                CASE WHEN (value -> 'Kind')::int = 0 THEN
                    jsonb_set(value, '{Kind}', '1')
                WHEN (value -> 'Kind')::int = 1 THEN
                    jsonb_set(value, '{Kind}', '2')
                WHEN (value -> 'Kind')::int = 2 THEN
                    jsonb_set(value, '{Kind}', '3')
                WHEN (value -> 'Kind')::int = 3 THEN
                    jsonb_set(value, '{Kind}', '4')
                WHEN (value -> 'Kind')::int = 4 THEN
                    jsonb_set(value, '{Kind}', '5')
                WHEN (value -> 'Kind')::int = 5 THEN
                    jsonb_set(value, '{Kind}', '6')
                WHEN (value -> 'Kind')::int = 6 THEN
                    jsonb_set(value, '{Kind}', '0')
                ELSE
                    value
                END)
        FROM
            jsonb_array_elements(enonce));

-- Put the constraint back
ALTER TABLE questions
    ADD CONSTRAINT enonce_structgen_validate_json_array_exe_Block CHECK (structgen_validate_json_array_exe_Block (enonce));

