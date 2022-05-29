-- v 0.4.0
-- require teacher/gen_create.sql

INSERT INTO teachers (id, mail, password_crypted, is_admin)
    VALUES (1, 'dummy@gmail.com', '\x5cc8e54c4df92cf6d49a86d85eff01b6b60e6e82d0db551fbb19295ebfd9bd15', TRUE);

SELECT
    setval('teachers_id_seq', (
            SELECT
                MAX(id)
            FROM teachers));

-- add teacher and visibility to trivial configs
-- visibility

ALTER TABLE trivial_configs
    ADD COLUMN public boolean;

UPDATE
    trivial_configs
SET
    public = TRUE;

ALTER TABLE trivial_configs
    ALTER COLUMN public SET NOT NULL;

-- admin : all the existing trivial_configs are owned by the admin account
ALTER TABLE trivial_configs
    ADD COLUMN id_teacher integer;

UPDATE
    trivial_configs
SET
    id_teacher = 1;

ALTER TABLE trivial_configs
    ALTER COLUMN id_teacher SET NOT NULL;

ALTER TABLE trivial_configs
    ADD FOREIGN KEY (id_teacher) REFERENCES teachers;

