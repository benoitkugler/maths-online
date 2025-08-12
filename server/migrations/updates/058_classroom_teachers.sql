BEGIN;
CREATE TABLE teacher_classrooms (
    IdTeacher integer NOT NULL,
    IdClassroom integer NOT NULL
);
ALTER TABLE teacher_classrooms
    ADD UNIQUE (IdTeacher, IdClassroom);
ALTER TABLE teacher_classrooms
    ADD FOREIGN KEY (IdTeacher) REFERENCES teachers;
ALTER TABLE teacher_classrooms
    ADD FOREIGN KEY (IdClassroom) REFERENCES classrooms;
INSERT INTO teacher_classrooms (IdClassroom, IdTeacher)
SELECT
    Id,
    IdTeacher
FROM
    classrooms;
-- related constraints
ALTER TABLE selfaccess_trivials
    DROP CONSTRAINT selfaccess_trivials_idclassroom_idteacher_fkey;
ALTER TABLE selfaccess_trivials
    ADD FOREIGN KEY (IdClassroom, IdTeacher) REFERENCES teacher_classrooms (IdClassroom, IdTeacher) ON DELETE CASCADE;
ALTER TABLE classrooms
    DROP COLUMN IdTeacher;
COMMIT;

