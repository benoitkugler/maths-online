package reviews

import (
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/homework"
	re "github.com/benoitkugler/maths-online/server/src/sql/reviews"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

type sample struct {
	adminID  uID
	userID   uID // own the items
	question editor.Questiongroup
	exercice editor.Exercicegroup
	trivial  trivial.Trivial
	sheet    homework.Sheet
}

var envs map[string]string

func init() {
	envs = tu.ReadEnv("../../../.env")
}

func setupDB(t *testing.T) (tu.TestDB, sample) {
	for k, v := range envs {
		t.Setenv(k, v)
	}

	userMail := os.Getenv("TEST_MAIL")
	tu.Assert(t, userMail != "")

	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql",
		"../../sql/trivial/gen_create.sql", "../../sql/homework/gen_create.sql", "../../sql/reviews/gen_create.sql")

	admin, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques, Mail: " "}.Insert(db)
	tu.AssertNoErr(t, err)
	user, err := teacher.Teacher{IsAdmin: false, FavoriteMatiere: teacher.Mathematiques, Mail: userMail}.Insert(db)
	tu.AssertNoErr(t, err)

	qu, err := editor.Questiongroup{IdTeacher: user.Id, Title: "Intervertion série intégrale (TEST)"}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, err := editor.Exercicegroup{IdTeacher: user.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	tr, err := trivial.Trivial{IdTeacher: user.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	sheet, err := homework.Sheet{IdTeacher: user.Id, Matiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	return db, sample{
		adminID:  admin.Id,
		userID:   user.Id,
		question: qu,
		exercice: ex,
		trivial:  tr,
		sheet:    sheet,
	}
}

func TestCRUDReviews(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r1, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.AssertNoErr(t, err)
	r2, err := ct.createReview(ReviewCreateIn{Kind: re.KExercice, Id: int64(sample.exercice.Id)}, sample.userID)
	tu.AssertNoErr(t, err)
	r3, err := ct.createReview(ReviewCreateIn{Kind: re.KTrivial, Id: int64(sample.trivial.Id)}, sample.userID)
	tu.AssertNoErr(t, err)
	r4, err := ct.createReview(ReviewCreateIn{Kind: re.KSheet, Id: int64(sample.sheet.Id)}, sample.userID)
	tu.AssertNoErr(t, err)

	ls, err := ct.listReviews()
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(ls) == 4)

	targets, err := re.LoadTargets(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(targets) == 4)

	err = ct.deleteReview(r1.Id, sample.userID)
	tu.AssertNoErr(t, err)

	err = ct.deleteReview(r2.Id, sample.adminID)
	tu.AssertNoErr(t, err)
	targets, err = re.LoadTargets(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(targets) == 2)

	err = ct.deleteReview(r3.Id, sample.userID)
	tu.AssertNoErr(t, err)
	err = ct.deleteReview(r4.Id, sample.userID)
	tu.AssertNoErr(t, err)
}

func TestCRUDReview(t *testing.T) {
	db, sample := setupDB(t)
	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.AssertNoErr(t, err)

	rExt, err := ct.loadReview(r.Id, sample.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(rExt.Comments) == 0)
	tu.Assert(t, rExt.Approvals == [3]int{})
	tu.Assert(t, rExt.IsAcceptable == false)
	tu.Assert(t, rExt.IsDeletable == true)

	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now(), Message: "Un premier message"},
		{Time: time.Now().Add(time.Hour), Message: "Un deuxième message"},
	}}, sample.adminID)
	tu.AssertNoErr(t, err)
	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now().Add(time.Minute), Message: "Une réponse"},
	}}, sample.userID)
	tu.AssertNoErr(t, err)

	err = ct.updateApproval(ReviewUpdateApprovalIn{IdReview: r.Id, Approval: re.Opposed}, sample.adminID)
	tu.AssertNoErr(t, err)
	err = ct.updateApproval(ReviewUpdateApprovalIn{IdReview: r.Id, Approval: re.InFavor}, sample.userID)
	tu.AssertNoErr(t, err)

	rExt, err = ct.loadReview(r.Id, sample.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(rExt.Comments) == 3)
	tu.Assert(t, rExt.Approvals == [3]int{0, 1, 1})
}

func TestAcceptReview(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()

	smtp, err := pass.NewSMTP()
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, smtp)

	r, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.AssertNoErr(t, err)

	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now(), Message: "Un premier message"},
		{Time: time.Now().Add(time.Hour), Message: "Un deuxième message"},
	}}, sample.adminID)
	tu.AssertNoErr(t, err)

	err = ct.acceptReview(r.Id, ct.admin.Id)
	tu.AssertNoErr(t, err)

	question, err := editor.SelectQuestiongroup(db, sample.question.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, question.Public)
	tu.Assert(t, question.IdTeacher == sample.adminID)
}

func TestLoadTarget(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()

	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r1, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.AssertNoErr(t, err)

	out1, err := ct.loadTargetContent(r1.Id, sample.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out1.(TargetQuestion).Group.Origin.Visibility == tcAPI.Personnal)

	out1, err = ct.loadTargetContent(r1.Id, sample.userID+10)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out1.(TargetQuestion).Group.Origin.Visibility == tcAPI.Hidden)

	err = ct.deleteReview(r1.Id, sample.userID)
	tu.AssertNoErr(t, err)

	r2, err := ct.createReview(ReviewCreateIn{Kind: re.KExercice, Id: int64(sample.exercice.Id)}, sample.userID)
	tu.AssertNoErr(t, err)

	out2, err := ct.loadTargetContent(r2.Id, sample.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out2.(TargetExercice).Group.Origin.Visibility == tcAPI.Personnal)

	out2, err = ct.loadTargetContent(r2.Id, sample.userID+10)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out2.(TargetExercice).Group.Origin.Visibility == tcAPI.Hidden)

	err = ct.deleteReview(r2.Id, sample.userID)
	tu.AssertNoErr(t, err)
}
