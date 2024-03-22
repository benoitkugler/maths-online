package teacher

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	evs "github.com/benoitkugler/maths-online/server/src/sql/events"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// This file defines the endpoints called by the student client app

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

	events, err := evs.SelectEventsByIdStudents(ct.db, student.Id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// display the teacher coordinates
	mail, url := teacher.Contact.Name, teacher.Contact.URL
	if mail == "" {
		mail, url = teacher.Mail, ""
	}
	return CheckStudentClassroomOut{
		IsOK: true,
		Meta: StudentClassroomHeader{
			Student: StudentClient{
				Name:    student.Name,
				Surname: student.Surname,
			},
			ClassroomName:     classroom.Name,
			TeacherMail:       mail,
			TeacherContactURL: url,
		},
		Advance: evs.NewAdvance(events).Stats(defaultMaxRankTreshold),
	}, nil
}

// AttachStudentToClassroomStep1 uses a temporary classroom code to
// attach a student to the classroom.
// More precisely, it checks the given code and returns a list of student
// propositions.
// As a special case, it also accepts a special demo code <DEMO_CODE>.[0-9] which creates a
// profile linked to the demo classroom.
func (ct *Controller) AttachStudentToClassroomStep1(c echo.Context) error {
	code := c.QueryParam("code")

	out, err := ct.attachStudentCandidates(code)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) checkClassroomCode(code string) (tc.IdClassroom, error) {
	item, ok, err := tc.SelectClassroomCodeByCode(ct.db, code)
	if err != nil {
		return 0, utils.SQLError(err)
	}
	if !ok {
		return 0, fmt.Errorf("Le code %s est invalide ou a expiré.", code)
	}
	return item.IdClassroom, nil
}

func (ct *Controller) attachStudentCandidates(code string) (AttachStudentToClassroom1Out, error) {
	// look for demonstration code
	if isDemoCode(ct.demoCode, code) {
		return ct.createDemoStudent()
	}

	idClassroom, err := ct.checkClassroomCode(code)
	if err != nil {
		return nil, err
	}

	stds, err := tc.SelectStudentsByIdClassrooms(ct.db, idClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// return the list of the students, indicating the ones already attached
	var out AttachStudentToClassroom1Out
	for _, student := range stds {
		out = append(out, NewStudentHeader(student))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out, nil
}

func isDemoCode(demoCode string, userCode string) bool {
	chunks := strings.Split(userCode, ".")
	if len(chunks) == 2 {
		return chunks[0] == demoCode
	}
	return false
}

func (ct *Controller) createDemoStudent() (AttachStudentToClassroom1Out, error) {
	student, err := tc.Student{
		Name:        "DEMO",
		Surname:     fmt.Sprintf("User %d", time.Now().Unix()),
		Birthday:    tc.Date(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)),
		IdClassroom: ct.demoClassroom.Id,
	}.Insert(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return AttachStudentToClassroom1Out{NewStudentHeader(student)}, nil
}

// AttachStudentToClassroomStep2 validates the birthday and actually attaches the client to
// a student account and a classroom.
func (ct *Controller) AttachStudentToClassroomStep2(c echo.Context) error {
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

func (ct *Controller) validAttachStudent(args AttachStudentToClassroom2In) (out AttachStudentToClassroom2Out, err error) {
	// check for expired codes
	if !isDemoCode(ct.demoCode, args.ClassroomCode) {
		_, err = ct.checkClassroomCode(args.ClassroomCode)
		if err != nil {
			return out, err
		}
	}

	student, err := tc.SelectStudent(ct.db, args.IdStudent)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// check if the birthday is correct
	if args.Birthday != time.Time(student.Birthday).Format(tc.DateLayout) {
		return AttachStudentToClassroom2Out{ErrInvalidBirthday: true}, nil
	}

	out = AttachStudentToClassroom2Out{
		IdCrypted: string(ct.studentKey.EncryptID(int64(args.IdStudent))),
	}

	student.Clients = append(student.Clients, tc.Client{
		Device: args.Device,
		Time:   time.Now(),
	})
	_, err = student.Update(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return out, nil
}

// StudentUpdatePlaylist is called to register an event when
// a student updates its playlist
func (ct *Controller) StudentUpdatePlaylist(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))
	id, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}

	notif, err := evs.RegisterEvents(ct.db, tc.IdStudent(id), evs.E_Misc_SetPlaylist)
	if err != nil {
		return err
	}

	return c.JSON(200, notif)
}

// in the future this might be overriden by the teacher
const defaultMaxRankTreshold = 40000
