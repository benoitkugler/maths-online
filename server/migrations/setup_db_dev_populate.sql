INSERT INTO questions (id, title, parameters, enonce)
    VALUES (1, 'Exemple 1', '{}', '[]'), (2, 'Exemple 2', '{}', '[]'), (3, 'Exemple 3', '{}', '[]'), (4, 'Exemple 4', '{}', '[]'), (5, 'Exemple 5', '{}', '[]'), (6, 'Exemple 6', '{}', '[]');

INSERT INTO question_tags (id_question, tag)
    VALUES (1, 'CAT1'), (1, 'CAT2'), (2, 'CAT1'), (2, 'CAT3'), (3, 'CAT3'), (4, 'CAT3'), (5, 'CAT3'), (6, 'CAT3');

INSERT INTO students (id, name, surname)
    VALUES (1, 'K', 'Benoit'), (2, 'L', 'GuiGui');

INSERT INTO trivial_configs (Id, Questions, QuestionTimeout)
    VALUES (1, '[[], [], [], [], []]', 60);

SELECT
    setval('questions_id_seq', (
            SELECT
                MAX(id)
            FROM questions));

SELECT
    setval('trivial_configs_id_seq', (
            SELECT
                MAX(id)
            FROM trivial_configs));

