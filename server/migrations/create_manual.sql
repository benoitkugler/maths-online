-- instructions to run once at DB creation
--  1) Admin teacher account
--  2) Demo classroom on the admin account

INSERT INTO teachers (id, mail, passwordCrypted, isAdmin, HasSimplifiedEditor, Contact, FavoriteMatiere)
    VALUES (1, 'XXX', 'XXX', TRUE, FALSE, '{"Name":"", "URL":""}', 'MATHS');

INSERT INTO classrooms (id, idTeacher, name, MaxRankThreshold)
    VALUES (1, 1, 'DEMO', 40000);

