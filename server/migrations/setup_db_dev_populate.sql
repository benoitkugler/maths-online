INSERT INTO teachers (id, mail, password_crypted, is_admin)
    VALUES (1, 'dummy@gmail.com', '\x5cc8e54c4df92cf6d49a86d85eff01b6b60e6e82d0db551fbb19295ebfd9bd15', TRUE);

INSERT INTO questions (id, page, public, id_teacher)
    VALUES (1, '{"title":"","enonce":[], "parameters":{}}', TRUE, 1);

INSERT INTO question_tags (id_question, tag)
    VALUES (1, 'CAT1'), (1, 'CAT2');

INSERT INTO students (id, name, surname)
    VALUES (1, 'K', 'Benoit'), (2, 'L', 'GuiGui');

INSERT INTO trivial_configs (Id, Questions, QuestionTimeout, ShowDecrassage)
    VALUES (1, '[[], [], [], [], []]', 60, TRUE);

SELECT
    setval('teachers_id_seq', (
            SELECT
                MAX(id)
            FROM teachers));

SELECT
    setval('questions_id_seq', (
            SELECT
                MAX(id)
            FROM questions));

SELECT
    setval('students_id_seq', (
            SELECT
                MAX(id)
            FROM students));

SELECT
    setval('trivial_configs_id_seq', (
            SELECT
                MAX(id)
            FROM trivial_configs));

