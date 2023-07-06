-- v1.5.1
BEGIN;
ALTER TABLE travails
    ADD UNIQUE (Id, IdSheet);
ALTER TABLE sheets
    ADD COLUMN Anonymous integer;
ALTER TABLE sheets
    ADD FOREIGN KEY (Anonymous) REFERENCES travails ON DELETE CASCADE;
ALTER TABLE sheets
    ADD FOREIGN KEY (Id, Anonymous) REFERENCES travails (IdSheet, Id) ON DELETE CASCADE;
COMMIT;

