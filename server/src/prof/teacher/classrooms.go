package teacher

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/benoitkugler/maths-online/pass"
	tc "github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

type ClassroomExt struct {
	Classroom tc.Classroom

	NbStudents int
}

func (ct *Controller) checkAcces(userID tc.IdTeacher, classroomID tc.IdClassroom) error {
	// check the access
	classroom, err := tc.SelectClassroom(ct.db, classroomID)
	if err != nil {
		return utils.SQLError(err)
	}

	if classroom.IdTeacher != userID {
		return accessForbidden
	}

	return nil
}

// check that the user has the ownership on the student
func (ct *Controller) checkStudentOwnership(userID tc.IdTeacher, studentID tc.IdStudent) error {
	student, err := tc.SelectStudent(ct.db, studentID)
	if err != nil {
		return utils.SQLError(err)
	}

	if err := ct.checkAcces(userID, student.IdClassroom); err != nil {
		return err
	}

	return nil
}

func (ct *Controller) TeacherGetClassrooms(c echo.Context) error {
	user := JWTTeacher(c)

	classrooms, err := tc.SelectClassroomsByIdTeachers(ct.db, user.Id)
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
	user := JWTTeacher(c)

	_, err := tc.Classroom{IdTeacher: user.Id, Name: "Nouvelle classe"}.Insert(ct.db)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) TeacherUpdateClassroom(c echo.Context) error {
	user := JWTTeacher(c)

	var args tc.Classroom
	if err := c.Bind(&args); err != nil {
		return err
	}

	// check the access
	if err := ct.checkAcces(user.Id, args.Id); err != nil {
		return err
	}

	args, err := args.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, args)
}

// TeacherDeleteClassroom remove the classrooms and all related students
func (ct *Controller) TeacherDeleteClassroom(c echo.Context) error {
	user := JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteClassroom(tc.IdClassroom(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteClassroom(idClassroom tc.IdClassroom, userID tc.IdTeacher) error {
	// check the access
	if err := ct.checkAcces(userID, idClassroom); err != nil {
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
	user := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	// check the access
	if err = ct.checkAcces(user.Id, tc.IdClassroom(idClassroom)); err != nil {
		return err
	}

	out, err := ct.getClassroomStudents(tc.IdClassroom(idClassroom))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getClassroomStudents(idClassroom tc.IdClassroom) ([]tc.Student, error) {
	stds, err := tc.SelectStudentsByIdClassrooms(ct.db, idClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make([]tc.Student, 0, len(stds))
	for _, student := range stds {
		out = append(out, student)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Surname < out[j].Surname })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })

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
	user := JWTTeacher(c)

	idClassroomS := c.FormValue("id-classroom")
	idClassroom, err := strconv.ParseInt(idClassroomS, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID parameter %s : %s", idClassroomS, err)
	}

	if err = ct.checkAcces(user.Id, tc.IdClassroom(idClassroom)); err != nil {
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
	user := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	out, err := ct.addStudent(tc.IdClassroom(idClassroom), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) addStudent(idClassroom tc.IdClassroom, userID tc.IdTeacher) (tc.Student, error) {
	// check the access
	if err := ct.checkAcces(userID, idClassroom); err != nil {
		return tc.Student{}, err
	}

	st, err := tc.Student{IdClassroom: idClassroom, Name: "Nouvel", Surname: "Eleve", Birthday: tc.Date(time.Now())}.Insert(ct.db)
	if err != nil {
		return tc.Student{}, utils.SQLError(err)
	}

	return st, nil
}

// TeacherDeleteStudent removes the student from the classroom and
// completely deletes it.
func (ct *Controller) TeacherDeleteStudent(c echo.Context) error {
	user := JWTTeacher(c)

	idStudent, err := utils.QueryParamInt64(c, "id-student")
	if err != nil {
		return err
	}

	err = ct.deleteStudent(tc.IdStudent(idStudent), user.Id)
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
	user := JWTTeacher(c)

	var args tc.Student
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateStudent(args, user.Id)
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

type classroomsCode struct {
	codes map[string]tc.IdClassroom // code for student -> id_classroom
	lock  sync.Mutex
}

func (cc *classroomsCode) newCode(idClassroom tc.IdClassroom) string {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	// generated the code
	code := utils.RandomID(true, 4, func(s string) bool {
		_, has := cc.codes[s]
		return has
	})
	// register it
	cc.codes[code] = idClassroom

	// time its removal
	time.AfterFunc(6*time.Hour, func() { cc.expireCode(code) })

	return code
}

func (cc *classroomsCode) expireCode(code string) {
	cc.lock.Lock()
	defer cc.lock.Unlock()
	delete(cc.codes, code)
}

// return the ID of the classroom
func (cc *classroomsCode) checkCode(code string) (tc.IdClassroom, error) {
	cc.lock.Lock()
	defer cc.lock.Unlock()

	out, ok := cc.codes[code]
	if !ok {
		return 0, fmt.Errorf("Le code %s est invalide ou a expiré.", code)
	}
	return out, nil
}

type GenerateClassroomCodeOut struct {
	Code string
}

// TeacherGenerateClassroomCode generates a temporary code to link students app
// to the classroom.
func (ct *Controller) TeacherGenerateClassroomCode(c echo.Context) error {
	user := JWTTeacher(c)

	idClassroom, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	// check the access
	if err = ct.checkAcces(user.Id, tc.IdClassroom(idClassroom)); err != nil {
		return err
	}

	code := ct.classCodes.newCode(tc.IdClassroom(idClassroom))
	out := GenerateClassroomCodeOut{Code: code}

	return c.JSON(200, out)
}

// ------------------------- student client API -------------------------

// CheckStudentClassroom is called on app startup, to check that the
// student credentials are still up to date.
func (ct *Controller) CheckStudentClassroom(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))

	out, err := ct.checkStudentClassroom(idCrypted)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) checkStudentClassroom(idCrypted pass.EncryptedID) (out CheckStudentClassroomOut, err error) {
	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		// maybe the ID is out of date
		return CheckStudentClassroomOut{IsOK: false}, nil
	}

	student, err := tc.SelectStudent(ct.db, tc.IdStudent(idStudent))
	if err == sql.ErrNoRows {
		// the student has been removed
		return CheckStudentClassroomOut{IsOK: false}, nil
	} else if err != nil {
		return out, utils.SQLError(err)
	}

	classroom, err := tc.SelectClassroom(ct.db, student.IdClassroom)
	if err != nil {
		return out, utils.SQLError(err)
	}

	teacher, err := tc.SelectTeacher(ct.db, classroom.IdTeacher)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return CheckStudentClassroomOut{
		IsOK: true,
		Meta: StudentClassroomHeader{
			Student:       student,
			ClassroomName: classroom.Name,
			TeacherMail:   teacher.Mail,
		},
	}, nil
}

// AttachStudentToClassroom1 uses a temporary classroom code to
// attach a student to the classroom.
func (ct *Controller) AttachStudentToClassroom1(c echo.Context) error {
	code := c.QueryParam("code")

	idClassroom, err := ct.classCodes.checkCode(code)
	if err != nil {
		return err
	}

	// return the list of the student who are not yet identified
	stds, err := tc.SelectStudentsByIdClassrooms(ct.db, idClassroom)
	if err != nil {
		return utils.SQLError(err)
	}

	var out AttachStudentToClassroom1Out
	for _, student := range stds {
		if student.IsClientAttached {
			continue
		}
		out = append(out, StudentHeader{Id: student.Id, Label: student.Name + " " + student.Surname})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })

	return c.JSON(200, out)
}

func (ct *Controller) validAttachStudent(args AttachStudentToClassroom2In) (out AttachStudentToClassroom2Out, err error) {
	_, err = ct.classCodes.checkCode(args.ClassroomCode)
	if err != nil {
		return out, err
	}

	student, err := tc.SelectStudent(ct.db, args.IdStudent)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// avoid usurpation
	if student.IsClientAttached {
		return AttachStudentToClassroom2Out{ErrAlreadyAttached: true}, nil
	}

	// check if the birthday is correct
	if args.Birthday != time.Time(student.Birthday).Format(tc.DateLayout) {
		return AttachStudentToClassroom2Out{ErrInvalidBirthday: true}, nil
	}

	out = AttachStudentToClassroom2Out{
		IdCrypted: string(ct.studentKey.EncryptID(int64(args.IdStudent))),
	}

	student.IsClientAttached = true
	_, err = student.Update(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return out, nil
}

// AttachStudentToClassroom2 validates the birthday and actually attaches the client to
// a student account and a classroom.
func (ct *Controller) AttachStudentToClassroom2(c echo.Context) error {
	var args AttachStudentToClassroom2In
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.validAttachStudent(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}
