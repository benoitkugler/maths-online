// Package events implements an event system
// to reward student progression across the whole app.
//
// See https://docs.google.com/spreadsheets/d/1nY7zKsZ6JjW51QDSbt7OgaVSqzL-UoBHu6Qp6C3P7XQ
package events

import (
	"database/sql"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type pointResolver interface {
	resolve(nbOccurences int) int
}

// all events have the same score
type constResolver int

func (cr constResolver) resolve(nbOccurences int) int {
	return int(cr) * nbOccurences
}

// only one event is taken into account
type uniqueResolver int

func (ur uniqueResolver) resolve(nbOccurences int) int {
	if nbOccurences >= 1 {
		return int(ur)
	}
	return 0
}

type eventProperties struct {
	basePoint int
	isUnique  bool
}

var properties = [NbEvents]pointResolver{
	E_IsyTriv_Create:       uniqueResolver(30),
	E_IsyTriv_Streak3:      constResolver(40),
	E_IsyTriv_Win:          constResolver(300),
	E_Homework_TaskDone:    constResolver(40),
	E_Homework_TravailDone: constResolver(100),
	E_All_QuestionRight:    constResolver(5),
	E_All_QuestionWrong:    constResolver(1),
	E_Misc_SetPlaylist:     uniqueResolver(30),
	E_ConnectStreak3:       constResolver(20),
	E_ConnectStreak7:       constResolver(50),
	E_ConnectStreak30:      constResolver(400),
}

var refT = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

func day(t time.Time) int {
	return int(t.Sub(refT).Hours() / 24)
}

type dayEvents struct {
	day    int // from the refT
	events []EventK
}

// isFlameDone returns true if the number of good answers
// is greater or equal to 3
func (de dayEvents) isFlameDone() bool {
	count := 0
	for _, event := range de.events {
		if event == E_All_QuestionRight {
			count++
		}
	}
	return count >= 3
}

// Advance is the stored list of events, for one student.
// Generally speaking, its features are dynamic : see the various access methods
type Advance []dayEvents // sorted by day

func NewAdvance(events Events) Advance {
	tmp := map[int][]EventK{}
	for _, event := range events {
		d := day(time.Time(event.Date))
		tmp[d] = append(tmp[d], event.Event)
	}

	out := make(Advance, 0, len(tmp))
	for day, l := range tmp {
		out = append(out, dayEvents{day, l})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].day < out[j].day })
	return out
}

// return the 3, 7 and 30 streaks (as a number of each)
// any kind of event counts
func (adv Advance) connectStreaks() (nb3, nb7, nb30 int) {
	if len(adv) == 0 {
		return
	}

	lastDay := adv[0].day
	streak3, streak7, streak30 := 1, 1, 1
	for _, day := range adv[1:] {
		if day.day != lastDay+1 { // not contiguous, reset the streaks
			streak3 = 1
			streak7 = 1
			streak30 = 1

			lastDay = day.day
			continue
		}

		// contiguous connection, increase the streaks...
		streak3++
		streak7++
		streak30++

		// ... do we have reached a threshold ?
		if streak3 == 3 {
			nb3++
			streak3 = 0
		}
		if streak7 == 7 {
			nb7++
			streak7 = 0
		}
		if streak30 == 30 {
			nb30++
			streak30 = 0
		}

		lastDay = day.day
	}

	return
}

// flames returns the number of consecutive days (containing the present day)
// for which (at least) 3 [E_All_QuestionRight] have been recorded
func (adv Advance) flames() int {
	start := day(time.Now())

	currentDay := start
	// loop backward from the last day
	for i := len(adv) - 1; i >= 0; i-- {
		item := adv[i]
		if item.day != currentDay { // no more contiguous
			break
		}
		if !item.isFlameDone() {
			break
		}

		// continue
		currentDay--
	}

	return start - currentDay
}

// for now, the occurences are sufficient to compute the points
func totalPoints(oc [NbEvents]int) int {
	total := 0
	for ev, nb := range oc {
		resolver := properties[ev]
		total += resolver.resolve(nb)
	}
	return total
}

// occurences returns the number of realisation for each event,
// including the dynamic ones (deduced from others).
func (adv Advance) occurences() (occurences [NbEvents]int) {
	for _, day := range adv {
		for _, ev := range day.events {
			occurences[ev]++
		}
	}

	// add the dynamic events
	nb3, nb7, nb30 := adv.connectStreaks()
	occurences[E_ConnectStreak3] = nb3
	occurences[E_ConnectStreak7] = nb7
	occurences[E_ConnectStreak30] = nb30

	return
}

type Stats struct {
	Occurences  [NbEvents]int
	TotalPoints int
	Flames      int
}

func (adv Advance) Stats() Stats {
	oc := adv.occurences()
	return Stats{
		Occurences:  oc,
		TotalPoints: totalPoints(oc),
		Flames:      adv.flames(),
	}
}

// RegisterEvents stores the given events for the given student at the present time.
//
// It returns the number of points earned by the student.
func RegisterEvents(db *sql.DB, idStudent teacher.IdStudent, events ...EventK) (EventNotification, error) {
	eventsBefore, err := SelectEventsByIdStudents(db, idStudent)
	if err != nil {
		return EventNotification{}, utils.SQLError(err)
	}
	advanceBefore := NewAdvance(eventsBefore)

	newEvents := make(Events, len(events))
	for i, ev := range events {
		newEvents[i] = Event{IdStudent: idStudent, Date: Date(time.Now()), Event: ev}
	}

	tx, err := db.Begin()
	if err != nil {
		return EventNotification{}, utils.SQLError(err)
	}
	err = InsertManyEvents(tx, newEvents...)
	if err != nil {
		_ = tx.Rollback()
		return EventNotification{}, utils.SQLError(err)
	}
	err = tx.Commit()
	if err != nil {
		return EventNotification{}, utils.SQLError(err)
	}

	eventsAfter := append(eventsBefore, newEvents...)
	advanceAfter := NewAdvance(eventsAfter)

	score := totalPoints(advanceAfter.occurences()) - totalPoints(advanceBefore.occurences())
	return EventNotification{Events: events, Points: score}, nil
}

type EventNotification struct {
	Events []EventK // the events (often with length 1)
	Points int      // the number of points earned
}

func (en *EventNotification) HideIfNoPoints() {
	if en.Points == 0 {
		en.Events = nil
	}
}
