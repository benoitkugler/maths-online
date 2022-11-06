package homework

import (
	"errors"
	"time"

	"github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/utils"
)

func LoadMonoquestionSheet(db DB, idMono tasks.IdMonoquestion) (tasks.IdTask, IdSheet, error) {
	ts, err := tasks.SelectTasksByIdMonoquestions(db, idMono)
	if err != nil {
		return 0, 0, utils.SQLError(err)
	}
	if len(ts) != 1 {
		return 0, 0, errors.New("internal error: expected one task for a monoquestion")
	}
	link, found, err := SelectSheetTaskByIdTask(db, ts.IDs()[0])
	if err != nil {
		return 0, 0, err
	}
	if !found {
		return 0, 0, errors.New("internal error: task without sheet")
	}
	return link.IdTask, link.IdSheet, nil
}

// IsExpired returns true if the [Deadline] is before the present time.
func (sh Sheet) IsExpired() bool {
	return time.Time(sh.Deadline).Before(time.Now())
}
