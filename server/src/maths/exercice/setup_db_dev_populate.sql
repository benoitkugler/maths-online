INSERT INTO questions (id, title, parameters, enonce)
    VALUES (1, 'Premier exemple', '{}', '[]'), (2, 'Deuxième exemple', '{}', '[]');

INSERT INTO question_tags (id_question, tag)
    VALUES (1, 'Seconde'), (1, 'Probas'), (2, 'Seconde'), (2, 'Calcul littéral');

SELECT
    setval('questions_id_seq', (
            SELECT
                MAX(id)
            FROM questions));

