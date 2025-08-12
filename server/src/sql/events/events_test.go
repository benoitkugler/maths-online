// Package events implements an event system
// to reward student progression across the whole app.
package events

import (
	"testing"
	"time"

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

func e(evs ...EventK) []EventK { return evs }

func TestAdvance_Flames(t *testing.T) {
	const E = E_All_QuestionRight
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
		if got := tt.adv.flames(); got != tt.want {
			t.Errorf("Advance.Flames() = %v, want %v", got, tt.want)
		}
	}
}

func TestSQL(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	cl, err := teacher.Classroom{}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)
	err = InsertManyEvents(tx, Events{
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Date(2023, time.September, 24, 1, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Date(2023, time.September, 24, 0, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Date(2023, time.September, 24, 23, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Date(2023, time.September, 25, 23, 2, 0, 0, time.Local))},
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Now())},
		{IdStudent: student.Id, Event: E_All_QuestionRight, Date: Date(time.Now())},
	}...)
	tu.AssertNoErr(t, err)
	err = tx.Commit()
	tu.AssertNoErr(t, err)

	l, err := SelectEventsByIdStudents(db, student.Id)
	tu.AssertNoErr(t, err)
	adv := NewAdvance(l)
	tu.Assert(t, len(adv) == 3)

	points, err := RegisterEvents(db.DB, student.Id, E_All_QuestionWrong, E_IsyTriv_Win)
	tu.AssertNoErr(t, err)
	tu.Assert(t, points.Points == 1+300)

	_, err = RegisterEvents(db.DB, student.Id, EventK(NbEvents+1))
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
		wantOccurences [NbEvents]int
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
			[11]int{0: 3, 1: 3, 2: 3, E_ConnectStreak3: 1},
		},
	}
	for _, tt := range tests {
		if gotOccurences := tt.adv.occurences(); gotOccurences != tt.wantOccurences {
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
		if got := totalPoints(tt.adv.occurences()); got != tt.want {
			t.Errorf("Advance.TotalPoints() = %v, want %v", got, tt.want)
		}
	}
}

func Test_progressionScale_findRank(t *testing.T) {
	tests := []struct {
		ratio float64
		want  int
	}{
		{0, 0},
		{4.99, 0},
		{5, 1},
		{15, 2},
		{30, 3},
		{50, 4},
		{73, 5},
		{99.999, 5},
		{100, 6},
		{110, 6},
	}
	// 0, 5, 15, 30, 50, 73, 100,
	for _, tt := range tests {
		if got := defaultProgressionScale.findRank(tt.ratio); got != tt.want {
			t.Errorf("progressionScale.findRank(%v) = %v, want %v", tt.ratio, got, tt.want)
		}
	}
}
