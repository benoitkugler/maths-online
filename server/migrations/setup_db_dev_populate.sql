INSERT INTO questions (id, title, parameters, enonce)
    VALUES (1, 'Exemple 1', '{}', '[]'), (2, 'Exemple 2', '{}', '[]'), (3, 'Exemple 3', '{}', '[]'), (4, 'Exemple 4', '{}', '[]'), (5, 'Exemple 5', '{}', '[]'), (6, 'Exemple 6', '{}', '[]');

INSERT INTO question_tags (id_question, tag)
    VALUES (1, 'Cat1'), (1, 'Cat2'), (2, 'Cat1'), (2, 'Cat3'), (3, 'Cat3'), (4, 'Cat3'), (5, 'Cat3'), (6, 'Cat3');

SELECT
    setval('questions_id_seq', (
            SELECT
                MAX(id)
            FROM questions));

