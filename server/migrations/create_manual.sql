-- instructions to run once at DB creation
--  1) Admin teacher account
--  2) Demo classroom on the admin account

INSERT INTO teachers (id, mail, passwordCrypted, isAdmin)
    VALUES (1, 'XXX', 'XXX', TRUE);

INSERT INTO classrooms (id, idTeacher, name)
    VALUES (1, 1, 'DEMO');

