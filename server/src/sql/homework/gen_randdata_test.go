package homework

import (
	"math/rand"
	"time"

	"github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/sql/teacher"
)

// Code generated by gomacro/generator/go/randdata. DO NOT EDIT.

func randIdSheet() IdSheet {
	return IdSheet(randint64())
}

func randNotation() Notation {
	choix := [...]Notation{NoNotation, SuccessNotation}
	i := rand.Intn(len(choix))
	return choix[i]
}

func randSheet() Sheet {
	return Sheet{
		Id:          randIdSheet(),
		IdClassroom: randtea_IdClassroom(),
		Title:       randstring(),
		Notation:    randNotation(),
		Activated:   randbool(),
		Deadline:    randTime(),
	}
}

func randSheetTask() SheetTask {
	return SheetTask{
		IdSheet: randIdSheet(),
		Index:   randint(),
		IdTask:  randtas_IdTask(),
	}
}

func randTime() Time {
	return Time(randtTime())
}

func randbool() bool {
	i := rand.Int31n(2)
	return i == 1
}

func randint() int {
	return int(rand.Intn(1000000))
}

func randint64() int64 {
	return int64(rand.Intn(1000000))
}

var letterRunes2 = []rune("azertyuiopqsdfghjklmwxcvbn123456789é@!?&èïab ")

func randstring() string {
	b := make([]rune, 50)
	maxLength := len(letterRunes2)
	for i := range b {
		b[i] = letterRunes2[rand.Intn(maxLength)]
	}
	return string(b)
}

func randtTime() time.Time {
	return time.Unix(int64(rand.Int31()), 5)
}

func randtas_IdTask() tasks.IdTask {
	return tasks.IdTask(randint64())
}

func randtea_IdClassroom() teacher.IdClassroom {
	return teacher.IdClassroom(randint64())
}