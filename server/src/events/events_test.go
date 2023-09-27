// Package events implements an event system
// to reward student progression across the whole app.
package events

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/events"
	evs "github.com/benoitkugler/maths-online/server/src/sql/events"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func mustT(s string) time.Time {
	out, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return out
}

func Test_day(t *testing.T) {
	tests := []struct {
		args time.Time
		want int
	}{
		{refT, 0},
		{refT.Add(time.Hour), 0},
		{refT.Add(23 * time.Hour), 0},
		{refT.Add(25 * time.Hour), 1},
		{mustT("2020-01-01"), 0},
		{mustT("2020-01-02"), 1},
	}
	for _, tt := range tests {
		if got := day(tt.args); got != tt.want {
			t.Errorf("day() = %v, want %v", got, tt.want)
		}
	}
}

func e(evs ...events.EventK) []events.EventK { return evs }

func TestAdvance_Flames(t *testing.T) {
	const E = events.E_All_QuestionRight
	present := day(time.Now())
	tests := []struct {
		adv  Advance
		want int
	}{
		{
			Advance{{present - 1, e(0, 0, E, E, E)}, {present, e(E, E, 0, E)}},
			2,
		},
		{
			Advance{{present - 1, e(0, 0, E, E, E, E)}, {present, e(E, E, 0, E, E)}},
			2,
		},
		{
			Advance{{present - 1, e(0, 0, E, E, E)}, {present - 1, e(E, E, E)}, {present, e(E, E, 0, 0)}},
			0,
		},
		{
			Advance{{present - 2, e(0, 0, E, E, E)}, {present - 1, e(E, E, E)}},
			0,
		},
		{
			Advance{{present - 1, e(0, 0, E, E, E)}, {present, e(E, E)}},
			0,
		},
		{
			Advance{{present - 1, e(0, 0, E, E)}, {present, e(E, E, E)}},
			1,
		},
	}
	for _, tt := range tests {
		if got := tt.adv.Flames(); got != tt.want {
			t.Errorf("Advance.Flames() = %v, want %v", got, tt.want)
		}
	}
}

func TestSQL(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/events/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{FavoriteMatiere: teacher.Allemand}.Insert(db)
	tu.AssertNoErr(t, err)

	cl, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)
	err = events.InsertManyEvents(tx, events.Events{
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Date(2023, time.September, 24, 1, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Date(2023, time.September, 24, 0, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Date(2023, time.September, 24, 23, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Date(2023, time.September, 25, 23, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Now())},
		{IdStudent: student.Id, Event: events.E_All_QuestionRight, Date: events.Date(time.Now())},
	}...)
	tu.AssertNoErr(t, err)
	err = tx.Commit()
	tu.AssertNoErr(t, err)

	l, err := events.SelectEventsByIdStudents(db, student.Id)
	tu.AssertNoErr(t, err)
	adv := NewAdvance(l)
	tu.Assert(t, len(adv) == 3)

	points, err := RegisterEvents(db.DB, student.Id, events.E_All_QuestionWrong, events.E_IsyTriv_Win)
	tu.AssertNoErr(t, err)
	tu.Assert(t, points == 1+300)

	_, err = RegisterEvents(db.DB, student.Id, events.EventK(events.NbEvents+1))
	tu.Assert(t, err != nil)
}

func TestAdvance_connectStreaks(t *testing.T) {
	cs := func(days ...int) Advance {
		var out Advance
		for _, d := range days {
			out = append(out, dayEvents{day: d, events: e(0)})
		}
		return out
	}
	tests := []struct {
		adv      Advance
		wantNb3  int
		wantNb7  int
		wantNb30 int
	}{
		{
			Advance{},
			0, 0, 0,
		},
		{
			cs(0),
			0, 0, 0,
		},
		{
			cs(0, 1, 2),
			1, 0, 0,
		},
		{
			cs(0, 1, 2, 3, 4, 5),
			2, 0, 0,
		},
		{
			cs(0, 1, 2, 13, 14, 15),
			2, 0, 0,
		},
		{
			cs(0, 1, 3, 4, 5, 6, 7, 8, 9),
			2, 1, 0,
		},
		{
			cs(0, 1, 2, 3, 4, 5, 6, 10, 11, 12, 13, 14, 15, 16),
			4, 2, 0,
		},
		{
			cs(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29),
			10, 4, 1,
		},
	}
	for _, tt := range tests {
		gotNb3, gotNb7, gotNb30 := tt.adv.connectStreaks()
		if gotNb3 != tt.wantNb3 {
			t.Errorf("Advance.connectStreaks() gotNb3 = %v, want %v", gotNb3, tt.wantNb3)
		}
		if gotNb7 != tt.wantNb7 {
			t.Errorf("Advance.connectStreaks() gotNb7 = %v, want %v", gotNb7, tt.wantNb7)
		}
		if gotNb30 != tt.wantNb30 {
			t.Errorf("Advance.connectStreaks() gotNb30 = %v, want %v", gotNb30, tt.wantNb30)
		}
	}
}

func TestAdvance_Events(t *testing.T) {
	tests := []struct {
		adv            Advance
		wantOccurences [evs.NbEvents]int
	}{
		{
			Advance{}, [11]int{},
		},
		{
			Advance{
				{0, e(0, 1, 2)},
			},
			[11]int{1, 1, 1},
		},
		{
			Advance{
				{0, e(0, 1, 2)},
				{1, e(0, 1, 2)},
			},
			[11]int{2, 2, 2},
		},
		{
			Advance{
				{0, e(0, 1, 2)},
				{1, e(0, 1, 2)},
				{2, e(0, 1, 2)},
			},
			[11]int{0: 3, 1: 3, 2: 3, events.E_ConnectStreak3: 1},
		},
	}
	for _, tt := range tests {
		if gotOccurences := tt.adv.Occurences(); gotOccurences != tt.wantOccurences {
			t.Errorf("Advance.Events() = %v, want %v", gotOccurences, tt.wantOccurences)
		}
	}
}

func TestAdvance_TotalPoints(t *testing.T) {
	tests := []struct {
		adv  Advance
		want int
	}{
		{
			Advance{}, 0,
		},
		{
			Advance{
				{0, e(1, 1, 1)},
			}, 3 * 40,
		},
		{
			Advance{
				{0, e(0, 0, 0)}, // unique
			}, 30,
		},
	}
	for _, tt := range tests {
		if got := tt.adv.TotalPoints(); got != tt.want {
			t.Errorf("Advance.TotalPoints() = %v, want %v", got, tt.want)
		}
	}
}
