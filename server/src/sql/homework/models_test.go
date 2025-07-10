package homework

import (
	"database/sql"
	"testing"
	"time"

	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func Test(t *testing.T) {
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	teacher, err := tc.Teacher{FavoriteMatiere: tc.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	cl, err := tc.Classroom{IdTeacher: teacher.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := tc.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	sheet, err := Sheet{IdTeacher: teacher.Id, Matiere: tc.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	travail, err := Travail{IdClassroom: cl.Id, IdSheet: sheet.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	item := TravailException{IdStudent: student.Id, IdTravail: travail.Id}
	err = item.Insert(db)
	tu.AssertNoErr(t, err)

	_, err = DeleteTravailExceptionsByIdStudentAndIdTravail(db, student.Id, travail.Id)
	tu.AssertNoErr(t, err)

	item = TravailException{IdStudent: student.Id, IdTravail: travail.Id, Deadline: sql.NullTime{
		Valid: true, Time: time.Now(),
	}}
	err = item.Insert(db)
	tu.AssertNoErr(t, err)

	item, ok, err := SelectTravailExceptionByIdStudentAndIdTravail(db, student.Id, travail.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, ok)
	tu.Assert(t, item.Deadline.Valid)
}
