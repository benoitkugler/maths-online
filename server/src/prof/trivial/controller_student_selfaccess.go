package trivial

import (
	"fmt"
	"sort"

	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// Student self access

type TrivialSelfaccess struct {
	Classrooms []teacher.Classroom
	Actives    []teacher.IdClassroom
}

// [TrivialGetSelfaccess] returns the classroom which may launch
// the given trivial
func (ct *Controller) TrivialGetSelfaccess(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id-trivial")
	if err != nil {
		return err
	}

	out, err := ct.selfaccess(trivial.IdTrivial(id), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) selfaccess(id trivial.IdTrivial, userID uID) (TrivialSelfaccess, error) {
	// load all classrooms for the user
	cls_, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return TrivialSelfaccess{}, utils.SQLError(err)
	}
	cls := make([]teacher.Classroom, 0, len(cls_))
	for _, cl := range cls_ {
		cls = append(cls, cl)
	}
	sort.Slice(cls, func(i, j int) bool { return cls[i].Name < cls[j].Name })

	links, err := trivial.SelectSelfaccessTrivialsByIdTrivials(ct.db, id)
	if err != nil {
		return TrivialSelfaccess{}, utils.SQLError(err)
	}

	actives := links.ByIdTeacher()[userID].IdClassrooms()
	return TrivialSelfaccess{Classrooms: cls, Actives: actives}, nil
}

type UpdateSelfaccessIn struct {
	IdTrivial    trivial.IdTrivial
	IdClassrooms []teacher.IdClassroom
}

// [TrivialUpdateSelfaccess] sets the classroom which may launch
// the given trivial
func (ct *Controller) TrivialUpdateSelfaccess(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var in UpdateSelfaccessIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	err := ct.updateSelfaccess(in, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSelfaccess(args UpdateSelfaccessIn, userID uID) error {
	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = trivial.DeleteSelfaccessTrivialsByIdTrivialAndIdTeacher(tx, args.IdTrivial, userID)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	newItems := make(trivial.SelfaccessTrivials, len(args.IdClassrooms))
	for i, id := range args.IdClassrooms {
		newItems[i] = trivial.SelfaccessTrivial{IdClassroom: id, IdTrivial: args.IdTrivial, IdTeacher: userID}
	}
	err = trivial.InsertManySelfaccessTrivials(tx, newItems...)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}
