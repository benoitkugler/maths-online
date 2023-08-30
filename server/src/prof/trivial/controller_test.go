package trivial

import (
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tr "github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestCreateConfig(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/trivial/gen_create.sql")

	tc, err := teacher.Teacher{FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	out, err := tr.Trivial{
		QuestionTimeout: 120,
		ShowDecrassage:  true,
		IdTeacher:       tc.Id,
	}.Insert(db)
	tu.AssertNoErr(t, err)

	if _, err := tr.DeleteTrivialById(db, out.Id); err != nil {
		t.Fatal(err)
	}
}

func TestGetConfig(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql",
		"../../sql/trivial/gen_create.sql", "../../sql/reviews/gen_create.sql")
	defer db.Remove()

	user1, err := teacher.Teacher{Mail: "1", FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)
	user2, err := teacher.Teacher{Mail: "2", FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	c1, err := tr.Trivial{IdTeacher: user1.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	_, err = tr.Trivial{IdTeacher: user2.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, pass.Encrypter{}, "", user1)
	l, err := ct.getTrivialPoursuits(user1.Id, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	if len(l) != 1 {
		t.Fatal(l)
	}

	c1.Public = true
	if _, err = c1.Update(db); err != nil {
		t.Fatal(err)
	}

	l, err = ct.getTrivialPoursuits(user2.Id, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	if len(l) != 2 {
		t.Fatal(l)
	}
}

func TestGameTermination(t *testing.T) {
	tv.ProgressLogger.SetOutput(os.Stdout)

	ct := newGameStore("test")

	ct.createGame(createGame{ID: selfaccessCode("Game1")})

	if len(ct.games) != 1 {
		t.Fatal("expected one game")
	}
	ct.stopGame(selfaccessCode("Game1"), false)

	time.Sleep(20 * time.Millisecond)

	ct.lock.Lock()
	if len(ct.games) != 0 {
		t.Fatal("game should have been removed")
	}
	ct.lock.Unlock()
}

func TestMissingQuestions(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{})

	criteria := tr.CategoriesQuestions{
		Tags: [...]tr.QuestionCriterion{
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "VIOLET"},
				},
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "VERT"},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "ORANGE"},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "JAUNE"},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "BLEU"},
				},
			},
		},
	}
	out, err := ct.checkMissingQuestions(criteria, 1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Missing) == 0)

	criteria = tr.CategoriesQuestions{
		Tags: [...]tr.QuestionCriterion{
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "VALEUR FINALE", Section: editor.TrivMath},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "TAUX RÉCIPROQUE", Section: editor.TrivMath},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "PROPORTION", Section: editor.TrivMath},
				},
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "PROPORTION DE PROPORTION", Section: editor.TrivMath},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "EVOLUTIONS IDENTIQUES", Section: editor.TrivMath},
				},
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "EVOLUTIONS SUCCESSIVES", Section: editor.TrivMath},
				},
			},
			{
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "COEFFICIENT MULTIPLICATEUR", Section: editor.TrivMath},
				},
				{
					{Tag: "POURCENTAGES", Section: editor.Chapter},
					{Tag: "TAUX D'ÉVOLUTION", Section: editor.TrivMath},
				},
			},
		},
	}
	out, err = ct.checkMissingQuestions(criteria, 1)
	tu.AssertNoErr(t, err)
	if len(out.Missing) == 0 {
		t.Fatal("categories should be missing")
	}
}

func TestGetTrivials(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{})

	for range [10]int{} {
		t.Run("", func(t *testing.T) {
			_, err := ct.getTrivialPoursuits(0, teacher.Mathematiques)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCRUDSelfaccess(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/trivial/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	tr1, err := tr.Trivial{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	tr2, err := tr.Trivial{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	cl1, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	cl2, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	cl3, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, pass.Encrypter{}, "", tc)

	l, err := ct.selfaccess(tr1.Id, tc.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Classrooms) == 3)
	tu.Assert(t, len(l.Actives) == 0)

	err = ct.updateSelfaccess(UpdateSelfaccessIn{IdTrivial: tr1.Id, IdClassrooms: []teacher.IdClassroom{cl1.Id}}, tc.Id)
	tu.AssertNoErr(t, err)

	err = ct.updateSelfaccess(UpdateSelfaccessIn{IdTrivial: tr2.Id, IdClassrooms: []teacher.IdClassroom{cl1.Id, cl3.Id}}, tc.Id)
	tu.AssertNoErr(t, err)

	err = ct.updateSelfaccess(UpdateSelfaccessIn{IdTrivial: tr1.Id, IdClassrooms: []teacher.IdClassroom{cl1.Id, cl2.Id, cl3.Id}}, tc.Id)
	tu.AssertNoErr(t, err)

	l, err = ct.selfaccess(tr1.Id, tc.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Classrooms) == 3)
	tu.Assert(t, len(l.Actives) == 3)

	l, err = ct.selfaccess(tr2.Id, tc.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Classrooms) == 3)
	tu.Assert(t, len(l.Actives) == 2)
}

func TestStudentSelfaccess(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/trivial/gen_create.sql", "../../sql/editor/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	tr1, err := tr.Trivial{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	tr2, err := tr.Trivial{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	_, err = tr.Trivial{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	cl1, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	st1, err := teacher.Student{IdClassroom: cl1.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, pass.Encrypter{}, "", tc)

	err = ct.updateSelfaccess(UpdateSelfaccessIn{IdTrivial: tr1.Id, IdClassrooms: []teacher.IdClassroom{cl1.Id}}, tc.Id)
	tu.AssertNoErr(t, err)
	err = ct.updateSelfaccess(UpdateSelfaccessIn{IdTrivial: tr2.Id, IdClassrooms: []teacher.IdClassroom{cl1.Id}}, tc.Id)
	tu.AssertNoErr(t, err)

	ls, err := ct.studentGetSelfaccess(st1.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(ls) == 2)

	out, err := ct.launchSelfaccess(tr1.Id, st1.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.GameID != "")
}
