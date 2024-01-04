package ceintures

import (
	"database/sql"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

type Controller struct {
	db    *sql.DB
	admin teacher.Teacher

	studentKey pass.Encrypter
}

func NewController(db *sql.DB, admin teacher.Teacher, studentKey pass.Encrypter) *Controller {
	return &Controller{
		db:         db,
		admin:      admin,
		studentKey: studentKey,
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

type GetSchemeOut struct {
	Scheme    Scheme
	Questions []ce.Beltquestion
}

func (ct *Controller) getScheme(userID uID) (GetSchemeOut, error) {
	// for now, there is only one scheme
	out := GetSchemeOut{Scheme: mathScheme}

	questions, err := ce.SelectAllBeltquestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	for _, qu := range questions {
		out.Questions = append(out.Questions, qu)
	}

	return out, nil
}

func (ct *Controller) CeinturesGetPending(c echo.Context) error {
	// userID := tcAPI.JWTTeacher(c)

	var args ce.Advance
	if err := c.Bind(&args); err != nil {
		return err
	}

	out := mathScheme.Pending(args)
	return c.JSON(200, out)
}
