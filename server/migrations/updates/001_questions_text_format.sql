-- replace IsHint by
-- Italic = true, Smaller = true

UPDATE
    questions
SET
    enonce = (
        SELECT
            jsonb_agg(
                CASE WHEN (value -> 'Kind')::int = 16 THEN
                    CASE WHEN (value -> 'Data' -> 'IsHint')::boolean = TRUE THEN
                        jsonb_set(value, '{Data}', (value -> 'Data') - 'IsHint' || '{"Italic": true, "Smaller":true}'::jsonb)
                    ELSE
                        jsonb_set(value, '{Data}', (value -> 'Data') - 'IsHint')
                    END
                ELSE
                    value
                END)
        FROM
            jsonb_array_elements(enonce));

