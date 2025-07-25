package teacher

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	evs "github.com/benoitkugler/maths-online/server/src/sql/events"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// 2 days
const classroomCodeDuration = 2 * 24 * time.Hour

type ClassroomExt struct {
	Classroom tc.Classroom

	NbStudents int
}

func (ct *Controller) checkAcces(userID tc.IdTeacher, classroomID tc.IdClassroom) (tc.Classroom, error) {
	// check the access
	classroom, err := tc.SelectClassroom(ct.db, classroomID)
	if err != nil {
		return tc.Classroom{}, utils.SQLError(err)
	}

	if classroom.IdTeacher != userID {
		return tc.Classroom{}, accessForbidden
	}

	return classroom, nil
}

// check that the user has the ownership on the student
func (ct *Controller) checkStudentOwnership(userID tc.IdTeacher, studentID tc.IdStudent) error {
	student, err := tc.SelectStudent(ct.db, studentID)
	if err != nil {
		return utils.SQLError(err)
	}

	if _, err := ct.checkAcces(userID, student.IdClassroom); err != nil {
		return err
	}

	return nil
}

func (ct *Controller) TeacherGetClassrooms(c echo.Context) error {
	userID := JWTTeacher(c)

	classrooms, err := tc.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return utils.SQLError(err)
	}

	students, err := tc.SelectStudentsByIdClassrooms(ct.db, classrooms.IDs()...)
	if err != nil {
		return utils.SQLError(err)
	}
	dict := students.ByIdClassroom()

	out := make([]ClassroomExt, 0, len(classrooms))
	for _, cl := range classrooms {
		out = append(out, ClassroomExt{Classroom: cl, NbStudents: len(dict[cl.Id])})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Classroom.Id < out[j].Classroom.Id })

	return c.JSON(200, out)
}

func (ct *Controller) TeacherCreateClassroom(c echo.Context) error {
	userID := JWTTeacher(c)

	_, err := tc.Classroom{
		IdTeacher:        userID,
		Name:             "Nouvelle classe",
		MaxRankThreshold: 40_000,
	}.Insert(ct.db)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) TeacherUpdateClassroom(c echo.Context) error {
	userID := JWTTeacher(c)

	var args tc.Classroom
	if err := c.Bind(&args); err != nil {
		return err
	}

	// check the access
	classroom, err := ct.checkAcces(userID, args.Id)
	if err != nil {
		return err
	}

	// basic check on MaxRank
	if !(100 <= args.MaxRankThreshold && args.MaxRankThreshold <= 1_000_000) {
		return errors.New("Seuil de la dernière guilde invalide")
	}

	classroom.Name = args.Name
	classroom.MaxRankThreshold = args.MaxRankThreshold

	classroom, err = classroom.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, classroom)
}

// TeacherDeleteClassroom remove the classrooms and all related students
func (ct *Controller) TeacherDeleteClassroom(c echo.Context) error {
	userID := JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteClassroom(tc.IdClassroom(id), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteClassroom(idClassroom tc.IdClassroom, userID tc.IdTeacher) error {
	// check the access
	if _, err := ct.checkAcces(userID, idClassroom); err != nil {
		return err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = tc.DeleteStudentsByIdClassrooms(tx, idClassroom)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	_, err = tc.DeleteClassroomById(tx, idClassroom)
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

func (ct *Controller) TeacherGetClassroomStudents(c echo.Context) error {
	userID := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	// check the access
	classroom, err := ct.checkAcces(userID, tc.IdClassroom(idClassroom))
	if err != nil {
		return err
	}

	out, err := ct.getClassroomStudents(classroom)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type StudentExt struct {
	Student tc.Student
	Success evs.StudentAdvance
}

func (ct *Controller) getClassroomStudents(classroom tc.Classroom) ([]StudentExt, error) {
	students, err := tc.SelectStudentsByIdClassrooms(ct.db, classroom.Id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	events, err := evs.SelectEventsByIdStudents(ct.db, students.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	m := events.ByIdStudent()

	out := make([]StudentExt, 0, len(students))
	for _, student := range students {
		out = append(out, StudentExt{
			Student: student,
			Success: evs.NewAdvance(m[student.Id]).Stats(classroom.MaxRankThreshold),
		})
	}

	// data is sorted client side, to handle new profiles

	return out, nil
}

// split NAME NAME Surname into NAME NAME ; Surname
// this format is used by the french Pronote service used in high schools
func parsePronoteName(s string) (name, surname string) {
	// name may be separated by spaces, ', -
	// so we prefer to rely on case
	chunks := strings.Fields(s)
	for i, chunk := range chunks {
		runes := []rune(chunk)
		if len(runes) < 2 {
			continue
		}
		if unicode.IsLetter(runes[1]) && unicode.IsLower(runes[1]) {
			// found the first surname
			name = strings.Join(chunks[:i], " ")
			surname = strings.Join(chunks[i:], " ")
			return
		}
	}

	// default to first chunk as name
	return chunks[0], strings.Join(chunks[1:], " ")
}

func parsePronoteStudentList(file io.Reader) ([]tc.Student, error) {
	const pronoteDateLayout = "02/01/2006"

	r := csv.NewReader(file)
	r.Comma = ';'

	lines, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Fichier d'élèves invalide : %s", err)
	}

	// remove the header
	if len(lines) < 1 {
		return nil, errors.New("Fichier d'élèves invalide : entête manquant.")
	}
	lines = lines[1:]

	out := make([]tc.Student, len(lines))
	for i, line := range lines {
		if len(line) < 2 {
			return nil, errors.New("Fichier d'élèves invalide : champs manquants")
		}
		name, surname := parsePronoteName(line[0])

		birthday, err := time.Parse(pronoteDateLayout, line[2])
		if err != nil {
			return nil, fmt.Errorf("Fichier d'élèves invalide (date) : %s", err)
		}

		out[i] = tc.Student{Name: name, Surname: surname, Birthday: tc.Date(birthday)}
	}

	return out, nil
}

// TeacherImportStudents import a CSV file generated by Pronote.
// Other formats could be added in the future.
func (ct *Controller) TeacherImportStudents(c echo.Context) error {
	userID := JWTTeacher(c)

	idClassroomS := c.FormValue("id-classroom")
	idClassroom, err := strconv.ParseInt(idClassroomS, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID parameter %s : %s", idClassroomS, err)
	}

	if _, err = ct.checkAcces(userID, tc.IdClassroom(idClassroom)); err != nil {
		return err
	}

	header, err := c.FormFile("file")
	if err != nil {
		return err
	}

	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	err = ct.importPronoteFile(file, tc.IdClassroom(idClassroom))
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) importPronoteFile(file io.Reader, idClassroom tc.IdClassroom) error {
	list, err := parsePronoteStudentList(file)
	if err != nil {
		return err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	for i, student := range list {
		student.IdClassroom = idClassroom

		list[i], err = student.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// TeacherAddStudent adds a new student to the given classroom.
func (ct *Controller) TeacherAddStudent(c echo.Context) error {
	userID := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	out, err := ct.addStudent(tc.IdClassroom(idClassroom), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) addStudent(idClassroom tc.IdClassroom, userID tc.IdTeacher) (StudentExt, error) {
	// check the access
	if _, err := ct.checkAcces(userID, idClassroom); err != nil {
		return StudentExt{}, err
	}

	st, err := tc.Student{IdClassroom: idClassroom, Name: "Nouvel", Surname: "Eleve", Birthday: tc.Date(time.Now())}.Insert(ct.db)
	if err != nil {
		return StudentExt{}, utils.SQLError(err)
	}

	return StudentExt{
		Student: st,
		Success: evs.StudentAdvance{},
	}, nil
}

// TeacherDeleteStudent removes the student from the classroom and
// completely deletes it.
func (ct *Controller) TeacherDeleteStudent(c echo.Context) error {
	userID := JWTTeacher(c)

	idStudent, err := utils.QueryParamInt64(c, "id-student")
	if err != nil {
		return err
	}

	err = ct.deleteStudent(tc.IdStudent(idStudent), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteStudent(idStudent tc.IdStudent, userID tc.IdTeacher) error {
	if err := ct.checkStudentOwnership(userID, idStudent); err != nil {
		return err
	}

	_, err := tc.DeleteStudentById(ct.db, idStudent)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// TeacherUpdateStudent updates the student profile.
func (ct *Controller) TeacherUpdateStudent(c echo.Context) error {
	userID := JWTTeacher(c)

	var args tc.Student
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateStudent(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateStudent(st tc.Student, userID tc.IdTeacher) error {
	if err := ct.checkStudentOwnership(userID, st.Id); err != nil {
		return err
	}

	// partial update
	existing, err := tc.SelectStudent(ct.db, st.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	existing.Name = st.Name
	existing.Surname = st.Surname
	existing.Birthday = st.Birthday
	_, err = existing.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

type GenerateClassroomCodeOut struct {
	Code string
}

// TeacherGenerateClassroomCode generates a temporary code to link students app
// to the classroom.
func (ct *Controller) TeacherGenerateClassroomCode(c echo.Context) error {
	userID := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	// check the access
	if _, err = ct.checkAcces(userID, tc.IdClassroom(idClassroom)); err != nil {
		return err
	}

	code, err := ct.generateClassroomCode(tc.IdClassroom(idClassroom))
	if err != nil {
		return err
	}
	out := GenerateClassroomCodeOut{Code: code}

	return c.JSON(200, out)
}

func (ct *Controller) generateClassroomCode(id tc.IdClassroom) (string, error) {
	// load the existing codes
	ccs, err := tc.SelectAllClassroomCodes(ct.db)
	if err != nil {
		return "", utils.SQLError(err)
	}
	m := ccs.Codes()

	// generate the code
	code := utils.RandomID(true, 4, func(s string) bool { return m[s] })

	// register it
	err = tc.ClassroomCode{
		IdClassroom: id,
		Code:        code,
		ExpiresAt:   tc.Time(time.Now().Add(classroomCodeDuration)),
	}.Insert(ct.db)
	if err != nil {
		return "", utils.SQLError(err)
	}

	// time its removal
	time.AfterFunc(classroomCodeDuration, func() { tc.CleanupClassroomCodes(ct.db) })

	return code, nil
}
