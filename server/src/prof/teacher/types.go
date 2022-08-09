package teacher

import (
	"time"
)

// Date represents a day, without time zone consideration
type Date time.Time

func (d Date) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Date) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

const DateLayout = "2006-01-02"

func (st Students) ByIdClassroom() map[IdClassroom]Students {
	out := make(map[IdClassroom]Students)
	for idStudent, student := range st {
		d := out[student.IdClassroom]
		if d == nil {
			d = make(Students)
		}
		d[idStudent] = student
		out[student.IdClassroom] = d
	}
	return out
}
