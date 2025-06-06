package ceintures

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/prof/preview"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

var errAccess = errors.New("access forbidden")

type Controller struct {
	db    *sql.DB
	admin teacher.Teacher

	studentKey pass.Encrypter

	anons anonymousEvolutions
}

func NewController(db *sql.DB, admin teacher.Teacher, studentKey pass.Encrypter) *Controller {
	return &Controller{
		db:         db,
		admin:      admin,
		studentKey: studentKey,
		anons:      newAnonymousEvolutions(),
	}
}

func (ct *Controller) CeinturesGetScheme(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	out, err := ct.getScheme(userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type BeltquestionHeader struct {
	Id    ce.IdBeltquestion
	Title string
}

type StageHeader struct {
	Questions []BeltquestionHeader
	HasTODO   bool
}

type GetSchemeOut struct {
	Scheme     Scheme
	Stages     [ce.NbDomains][ce.NbRanks]StageHeader
	IsAdmin    bool
	Classrooms []teacher.Classroom
}

func (ct *Controller) getScheme(userID uID) (GetSchemeOut, error) {
	// for now, there is only one scheme
	out := GetSchemeOut{Scheme: mathScheme, IsAdmin: userID == ct.admin.Id}

	classrooms, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	for _, cl := range classrooms {
		out.Classrooms = append(out.Classrooms, cl)
	}
	sort.Slice(out.Classrooms, func(i, j int) bool { return out.Classrooms[i].Id < out.Classrooms[j].Id })

	questions, err := ce.SelectAllBeltquestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	for stage, l := range byStage(questions) {
		var hasTODO bool
		qus := make([]BeltquestionHeader, len(l))
		for i, qu := range l {
			qus[i] = BeltquestionHeader{qu.Id, qu.Title}
			if qu.Parameters.HasTODO() {
				hasTODO = true
			}
		}
		out.Stages[stage.Domain][stage.Rank] = StageHeader{qus, hasTODO}
	}

	return out, nil
}

func (ct *Controller) CeinturesGetPending(c echo.Context) error {
	// userID := tcAPI.JWTTeacher(c)

	var args ce.Advance
	if err := c.Bind(&args); err != nil {
		return err
	}

	out := mathScheme.Pending(args, ce.Seconde)
	return c.JSON(200, out)
}

func (ct *Controller) CeinturesGetQuestions(c echo.Context) error {
	var args Stage
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.getQuestions(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// sorted by Id
func (ct *Controller) getQuestions(stage Stage) ([]ce.Beltquestion, error) {
	questions, err := ce.SelectBeltquestionsByDomainAndRank(ct.db, stage.Domain, stage.Rank)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	out := make([]ce.Beltquestion, 0, len(questions))
	for _, qu := range questions {
		out = append(out, qu)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })

	return out, nil
}

func (ct *Controller) CeinturesCreateQuestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	if userID != ct.admin.Id {
		return errAccess
	}

	var args Stage
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.createQuestion(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createQuestion(stage Stage) (ce.Beltquestion, error) {
	qu, err := ce.Beltquestion{
		Domain: stage.Domain,
		Rank:   stage.Rank,
		Repeat: 1,
	}.Insert(ct.db)
	if err != nil {
		return qu, utils.SQLError(err)
	}
	return qu, nil
}

type UpdateBeltquestionIn struct {
	Id     ce.IdBeltquestion
	Repeat int
	Title  string
}

func (ct *Controller) CeinturesUpdateQuestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	if userID != ct.admin.Id {
		return errAccess
	}

	var args UpdateBeltquestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	qu, err := ce.SelectBeltquestion(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	qu.Repeat = args.Repeat
	qu.Title = args.Title

	_, err = qu.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

// CeinturesDuplicateQuestion duplicates the given question into
// the same stage
func (ct *Controller) CeinturesDuplicateQuestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	if userID != ct.admin.Id {
		return errAccess
	}

	id, err := utils.QueryParamInt64(c, "id-question")
	if err != nil {
		return err
	}

	out, err := ct.duplicateQuestion(ce.IdBeltquestion(id))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type DuplicateQuestionOut struct {
	Question ce.Beltquestion
	Preview  preview.LoopbackShowCeinture
}

func (ct *Controller) duplicateQuestion(id ce.IdBeltquestion) (DuplicateQuestionOut, error) {
	question, err := ce.SelectBeltquestion(ct.db, id)
	if err != nil {
		return DuplicateQuestionOut{}, utils.SQLError(err)
	}
	out := question
	out, err = out.Insert(ct.db)
	if err != nil {
		return DuplicateQuestionOut{}, utils.SQLError(err)
	}

	preview, err := ct.preview(Stage{out.Domain, out.Rank}, false, out.Id)
	if err != nil {
		return DuplicateQuestionOut{}, err
	}

	return DuplicateQuestionOut{out, preview}, nil
}

type SaveBeltQuestionIn struct {
	Question       ce.Beltquestion
	ShowCorrection bool
}

func (ct *Controller) CeinturesSaveQuestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	if userID != ct.admin.Id {
		return errAccess
	}

	var args SaveBeltQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.saveQuestion(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type SaveBeltquestionAndPreviewOut struct {
	Error   questions.ErrQuestionInvalid
	IsValid bool
	Preview preview.LoopbackShowCeinture
}

func (ct *Controller) saveQuestion(args SaveBeltQuestionIn) (SaveBeltquestionAndPreviewOut, error) {
	qu := args.Question
	if err := qu.Page().Validate(); err != nil {
		return SaveBeltquestionAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// TODO: only preview for non admin members

	// save the question and load the group
	_, err := qu.Update(ct.db)
	if err != nil {
		return SaveBeltquestionAndPreviewOut{}, utils.SQLError(err)
	}

	pr, err := ct.preview(Stage{qu.Domain, qu.Rank}, args.ShowCorrection, qu.Id)
	if err != nil {
		return SaveBeltquestionAndPreviewOut{}, err
	}

	return SaveBeltquestionAndPreviewOut{IsValid: true, Preview: pr}, nil
}

func (ct *Controller) preview(stage Stage, showCorrection bool, currentQuestion ce.IdBeltquestion) (out preview.LoopbackShowCeinture, _ error) {
	l, err := ct.getQuestions(stage)
	if err != nil {
		return out, err
	}
	out.Origin = make([]questions.QuestionPage, len(l))
	for i, qu := range l {
		out.Origin[i] = qu.Page()
	}
	out.Questions, err = instantiateQuestions(l)
	if err != nil {
		return out, err
	}
	out.ShowCorrection = showCorrection
	// select the proper question
	for index, qu := range out.Questions {
		if qu.Id == currentQuestion {
			out.QuestionIndex = index
			break
		}
	}

	return out, nil
}

func (ct *Controller) LoopbackEvaluateCeinture(c echo.Context) error {
	var args preview.LoopbackEvaluateCeintureIn

	if err := c.Bind(&args); err != nil {
		return err
	}

	res, err := tasks.EvaluateBelt(ct.db, args.Questions, args.Answers)
	if err != nil {
		return err
	}

	out := preview.LoopbackEvaluateCeintureOut{Answers: res}
	return c.JSON(200, out)
}

// required to properly handle interface to JSON
type DeleteBeltquestionOut struct {
	Preview preview.LoopbackServerEvent
}

func (d DeleteBeltquestionOut) MarshalJSON() ([]byte, error) {
	data := struct {
		Preview preview.LoopbackServerEventWrapper
	}{preview.LoopbackServerEventWrapper{Data: d.Preview}}
	return json.Marshal(data)
}

func (ct *Controller) CeinturesDeleteQuestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	if userID != ct.admin.Id {
		return errAccess
	}

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.deleteQuestion(ce.IdBeltquestion(id))
	if err != nil {
		return err
	}

	return c.JSON(200, DeleteBeltquestionOut{out})
}

func (ct *Controller) deleteQuestion(id ce.IdBeltquestion) (preview.LoopbackServerEvent, error) {
	deleted, err := ce.DeleteBeltquestionById(ct.db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	stage := Stage{deleted.Domain, deleted.Rank}
	others, err := ct.getQuestions(stage)
	if err != nil {
		return nil, err
	}

	if len(others) == 0 { // nothing to show
		return preview.LoopbackPaused{}, nil
	}

	// select the first question
	preview, err := ct.preview(stage, false, others[0].Id)
	if err != nil {
		return nil, err
	}

	return preview, nil
}

// CeinturesGetStudentsAdvance returns the belt advance for the given classroom
func (ct *Controller) CeinturesGetStudentsAdvance(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id_, err := utils.QueryParamInt64(c, "classroom-id")
	if err != nil {
		return err
	}
	out, err := ct.getStudentsAdvance(teacher.IdClassroom(id_), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type studentAdvance struct {
	Student teacher.Student
	Advance ce.Beltevolution
}

func (ct *Controller) getStudentsAdvance(id teacher.IdClassroom, user uID) ([]studentAdvance, error) {
	classroom, err := teacher.SelectClassroom(ct.db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	if classroom.IdTeacher != user {
		return nil, errAccess
	}

	students, err := teacher.SelectStudentsByIdClassrooms(ct.db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links, err := ce.SelectBeltevolutionsByIdStudents(ct.db, students.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	byStudent := links.ByIdStudent()

	out := make([]studentAdvance, 0, len(students))
	for _, student := range students {
		evolution := byStudent[student.Id]
		out = append(out, studentAdvance{student, evolution})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Student.Name < out[j].Student.Name })
	return out, nil
}
