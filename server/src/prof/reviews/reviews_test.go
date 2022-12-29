package reviews

import (
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
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
		"../../sql/trivial/gen_create.sql", "../../sql/reviews/gen_create.sql")

	admin, err := teacher.Teacher{IsAdmin: true, Mail: " "}.Insert(db)
	tu.Assert(t, err == nil)
	user, err := teacher.Teacher{IsAdmin: false, Mail: userMail}.Insert(db)
	tu.Assert(t, err == nil)

	qu, err := editor.Questiongroup{IdTeacher: user.Id, Title: "Intervertion série intégrale"}.Insert(db)
	tu.Assert(t, err == nil)

	ex, err := editor.Exercicegroup{IdTeacher: user.Id}.Insert(db)
	tu.Assert(t, err == nil)

	tr, err := trivial.Trivial{IdTeacher: user.Id}.Insert(db)
	tu.Assert(t, err == nil)

	return db, sample{
		adminID:  admin.Id,
		userID:   user.Id,
		question: qu,
		exercice: ex,
		trivial:  tr,
	}
}

func TestCRUDReviews(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r1, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.Assert(t, err == nil)
	r2, err := ct.createReview(ReviewCreateIn{Kind: re.KExercice, Id: int64(sample.exercice.Id)}, sample.userID)
	tu.Assert(t, err == nil)
	r3, err := ct.createReview(ReviewCreateIn{Kind: re.KTrivial, Id: int64(sample.trivial.Id)}, sample.userID)
	tu.Assert(t, err == nil)

	ls, err := ct.listReviews()
	tu.Assert(t, err == nil)
	tu.Assert(t, len(ls) == 3)

	targets, err := re.LoadTargets(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(targets) == 3)

	err = ct.deleteReview(r1.Id, sample.userID)
	tu.Assert(t, err == nil)

	err = ct.deleteReview(r2.Id, sample.adminID)
	tu.Assert(t, err == nil)
	targets, err = re.LoadTargets(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(targets) == 1)

	err = ct.deleteReview(r3.Id, sample.userID)
	tu.Assert(t, err == nil)
}

func TestCRUDReview(t *testing.T) {
	db, sample := setupDB(t)
	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.Assert(t, err == nil)

	rExt, err := ct.loadReview(r.Id, sample.userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(rExt.Comments) == 0)
	tu.Assert(t, rExt.Approvals == [3]int{})
	tu.Assert(t, rExt.IsAcceptable == false)
	tu.Assert(t, rExt.IsDeletable == true)

	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now(), Message: "Un premier message"},
		{Time: time.Now().Add(time.Hour), Message: "Un deuxième message"},
	}}, sample.adminID)
	tu.Assert(t, err == nil)
	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now().Add(time.Minute), Message: "Une réponse"},
	}}, sample.userID)
	tu.Assert(t, err == nil)

	err = ct.updateApproval(ReviewUpdateApprovalIn{IdReview: r.Id, Approval: re.Opposed}, sample.adminID)
	tu.Assert(t, err == nil)
	err = ct.updateApproval(ReviewUpdateApprovalIn{IdReview: r.Id, Approval: re.InFavor}, sample.userID)
	tu.Assert(t, err == nil)

	rExt, err = ct.loadReview(r.Id, sample.userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(rExt.Comments) == 3)
	tu.Assert(t, rExt.Approvals == [3]int{0, 1, 1})
}

func TestAcceptReview(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()

	smtp, err := pass.NewSMTP()
	tu.Assert(t, err == nil)

	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, smtp)

	r, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.Assert(t, err == nil)

	err = ct.updateReview(ReviewUpdateCommentsIn{IdReview: r.Id, Comments: re.Comments{
		{Time: time.Now(), Message: "Un premier message"},
		{Time: time.Now().Add(time.Hour), Message: "Un deuxième message"},
	}}, sample.adminID)
	tu.Assert(t, err == nil)

	err = ct.acceptReview(r.Id, ct.admin.Id)
	tu.Assert(t, err == nil)

	question, err := editor.SelectQuestiongroup(db, sample.question.Id)
	tu.Assert(t, err == nil)
	tu.Assert(t, question.Public)
	tu.Assert(t, question.IdTeacher == sample.adminID)
}

func TestLoadTarget(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()

	ct := NewController(db.DB, teacher.Teacher{Id: sample.adminID}, pass.SMTP{})

	r1, err := ct.createReview(ReviewCreateIn{Kind: re.KQuestion, Id: int64(sample.question.Id)}, sample.userID)
	tu.Assert(t, err == nil)

	out1, err := ct.loadTargetContent(r1.Id, sample.userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, out1.(TargetQuestion).Group.Origin.Visibility == tcAPI.Personnal)

	out1, err = ct.loadTargetContent(r1.Id, sample.userID+10)
	tu.Assert(t, err == nil)
	tu.Assert(t, out1.(TargetQuestion).Group.Origin.Visibility == tcAPI.Hidden)

	err = ct.deleteReview(r1.Id, sample.userID)
	tu.Assert(t, err == nil)

	r2, err := ct.createReview(ReviewCreateIn{Kind: re.KExercice, Id: int64(sample.exercice.Id)}, sample.userID)
	tu.Assert(t, err == nil)

	out2, err := ct.loadTargetContent(r2.Id, sample.userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, out2.(TargetExercice).Group.Origin.Visibility == tcAPI.Personnal)

	out2, err = ct.loadTargetContent(r2.Id, sample.userID+10)
	tu.Assert(t, err == nil)
	tu.Assert(t, out2.(TargetExercice).Group.Origin.Visibility == tcAPI.Hidden)

	err = ct.deleteReview(r2.Id, sample.userID)
	tu.Assert(t, err == nil)
}
