-- v 0.4.1
ALTER TABLE trivial_configs
    ADD COLUMN Name varchar;

UPDATE
    trivial_configs
SET
    Name = 'Trivial ' || (id::text);

ALTER TABLE trivial_configs
    ALTER COLUMN Name SET NOT NULL;

