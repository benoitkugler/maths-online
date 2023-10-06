package homework

import (
	"errors"
	"fmt"
	"sort"

	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// this file provides various statistics for a [Travail],
// computed from the student progressions

type HowemorkMarksIn struct {
	IdClassroom tc.IdClassroom
	IdTravaux   []ho.IdTravail
}

type StudentTravailMark struct {
	Mark      float64 // /20
	Dispensed bool    // true if the student is dispensed for this travail
	NbTries   int     // the total number of tries (both success and failures) on this sheet
}

type TravailMarks struct {
	Marks     map[tc.IdStudent]StudentTravailMark // the notes for each student
	TaskStats []TaskStat
}

type TaskStat struct {
	IdWork taAPI.WorkID
	Title  string // title

	QuestionStats []QuestionStat

	// the total number of answers, for the whole "exercice"
	// it is the sum of each questions, provided for convenience
	NbSuccess, NbFailure int
}

func (ts *TaskStat) inferTotal() {
	for _, qu := range ts.QuestionStats {
		ts.NbSuccess += qu.NbSuccess
		ts.NbFailure += qu.NbFailure
	}
}

type QuestionStat struct {
	Description          string
	Id                   editor.IdQuestion
	Difficulty           editor.DifficultyTag
	NbSuccess, NbFailure int // the total number of answers, for this question
}

type HomeworkMarksOut struct {
	Students []tcAPI.StudentHeader // the students of the classroom
	Marks    map[ho.IdTravail]TravailMarks
}

func (ct *Controller) HomeworkGetMarks(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args HowemorkMarksIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.getMarks(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getMarks(args HowemorkMarksIn, userID uID) (HomeworkMarksOut, error) {
	classroom, err := tc.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}
	if classroom.IdTeacher != userID {
		return HomeworkMarksOut{}, errAccessForbidden
	}

	// returns the students for one classroom,
	// sorted alphabetically.
	stds, err := tc.SelectStudentsByIdClassrooms(ct.db, args.IdClassroom)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}

	travaux, err := ho.SelectTravails(ct.db, args.IdTravaux...)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}

	// load the dispenses
	links, err := ho.SelectTravailExceptionsByIdTravails(ct.db, travaux.IDs()...)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}
	expects := links.ByIdTravail()

	out := HomeworkMarksOut{
		Students: make([]tcAPI.StudentHeader, 0, len(stds)),
		Marks:    make(map[ho.IdTravail]TravailMarks),
	}
	// student list
	for _, s := range stds {
		out.Students = append(out.Students, tcAPI.NewStudentHeader(s))
	}

	sort.Slice(out.Students, func(i, j int) bool { return out.Students[i].Label < out.Students[j].Label })

	// compute the sheets marks :
	loader, err := newSheetsLoader(ct.db, travaux.IdSheets())
	if err != nil {
		return HomeworkMarksOut{}, err
	}
	// load all the progressions : for each task and student
	progressions, err := loader.tasks.LoadProgressions(ct.db)
	if err != nil {
		return HomeworkMarksOut{}, err
	}

	for idTravail, travail := range travaux {
		if travail.IdClassroom != classroom.Id {
			return HomeworkMarksOut{}, errors.New("internal error: inconsitent classroom ID")
		}

		markByStudent := make(map[tc.IdStudent]StudentTravailMark)
		var sheetTotal int
		// for each student, get its progression for each task
		tasks := loader.tasksForSheet(travail.IdSheet)

		tm := TravailMarks{
			Marks: markByStudent,
		}

		for _, link := range tasks {
			task := loader.tasks.Tasks[link.IdTask]
			work := loader.tasks.GetWork(task)
			bareme := work.Bareme()
			taskTotal := bareme.Total()
			sheetTotal += taskTotal
			byStudent := progressions[link.IdTask]

			questionsRes := make(map[editor.IdQuestion][2]int) // success, failure
			questions := make(editor.Questions)                // success, failure

			// add each progression to the student note
			for _, student := range stds { // make sure to consider all students
				studentProg := byStudent[student.Id]
				studentMark := bareme.ComputeMark(studentProg.Questions)
				item := markByStudent[student.Id]
				item.Mark += float64(studentMark)
				item.NbTries += studentProg.NbTries()
				markByStudent[student.Id] = item

				// map each question to its origin and compute its stats
				questionOrigins := loader.tasks.ResolveQuestions(student.Id, work)
				for questionIndex, origin := range questionOrigins {
					questions[origin.Id] = origin

					if len(studentProg.Questions) != 0 {
						succes, failure := studentProg.Questions[questionIndex].Stats()
						v := questionsRes[origin.Id]
						v[0] += succes
						v[1] += failure
						questionsRes[origin.Id] = v
					}
				}
			}

			taskStat := TaskStat{
				IdWork: taAPI.NewWorkID(task),
				Title:  work.Title(),
			}
			for index, qu := range loader.tasks.OrderQuestions(work) {
				res := questionsRes[qu.Id]
				stat := QuestionStat{
					Id:         qu.Id,
					Difficulty: qu.Difficulty,
					NbSuccess:  res[0], NbFailure: res[1],
				}
				switch taskStat.IdWork.Kind {
				case taAPI.WorkExercice:
					stat.Description = fmt.Sprintf("Question %d", index+1)
				case taAPI.WorkMonoquestion, taAPI.WorkRandomMonoquestion:
					stat.Description = qu.Subtitle
				}

				taskStat.QuestionStats = append(taskStat.QuestionStats, stat)
			}

			taskStat.inferTotal()
			tm.TaskStats = append(tm.TaskStats, taskStat)
		}

		// normalize the mark / 20 and add dispenses
		exceptions := expects[idTravail].ByIdStudent()
		for id, item := range markByStudent {
			item.Mark = 20 * item.Mark / float64(sheetTotal)
			if l := exceptions[id]; len(l) != 0 {
				item.Dispensed = l[0].IgnoreForMark
			}
			markByStudent[id] = item
		}

		out.Marks[idTravail] = tm
	}

	return out, nil
}
