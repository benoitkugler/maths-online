BEGIN;
-- bug fix : add key cascade
ALTER TABLE random_monoquestion_variants
    DROP CONSTRAINT random_monoquestion_variants_idstudent_fkey;
ALTER TABLE random_monoquestion_variants
    DROP CONSTRAINT random_monoquestion_variants_idrandommonoquestion_fkey;
ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdStudent) REFERENCES students ON DELETE CASCADE;
ALTER TABLE random_monoquestion_variants
    ADD FOREIGN KEY (IdRandomMonoquestion) REFERENCES random_monoquestions ON DELETE CASCADE;
-- optimize storage
ALTER TABLE progressions
    ALTER COLUMN INDEX TYPE smallint;
ALTER TABLE random_monoquestion_variants
    ALTER COLUMN INDEX TYPE smallint;
COMMIT;

