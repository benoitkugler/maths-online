// Package events implements an event system
// to reward student progression across the whole app.
package events

import (
	"sort"
	"time"

	evs "github.com/benoitkugler/maths-online/server/src/sql/events"
)

type eventProperties struct {
	basePoint int
	isUnique  bool
}

var properties = [...]eventProperties{
	evs.E_IsyTriv_Create:       {30, true},
	evs.E_IsyTriv_Streak3:      {40, false},
	evs.E_IsyTriv_Win:          {300, false},
	evs.E_Homework_TaskDone:    {40, false},
	evs.E_Homework_TravailDone: {100, false},
	evs.E_All_QuestionRight:    {5, false},
	evs.E_All_QuestionWrong:    {1, false},
	evs.E_Misc_SetPlaylist:     {30, true},
	evs.E_ConnectStreak3:       {20, false},
	evs.E_ConnectStreak7:       {50, false},
	evs.E_ConnectStreak30:      {400, false},
}

var refT = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

func day(t time.Time) int {
	return int(t.Sub(refT).Hours() / 24)
}

type dayEvents struct {
	day    int // from the refT
	events []evs.EventK
}

// isFlameDone returns true if the number of good answers
// is greater or equal to 3
func (de dayEvents) isFlameDone() bool {
	count := 0
	for _, event := range de.events {
		if event == evs.E_All_QuestionRight {
			count++
		}
	}
	return count >= 3
}

// Advance is the stored list of events, for one student.
// Generally speaking, its features are dynamic : see the various access methods
type Advance []dayEvents // sorted by day

func NewAdvance(events evs.Events) Advance {
	tmp := map[int][]evs.EventK{}
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

// Flames returns the number of consecutive days (containing the present day)
// for which (at least) 3 [E_All_QuestionRight] have been recorded
func (adv Advance) Flames() int {
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
