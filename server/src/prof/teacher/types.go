package teacher

import (
	"encoding/json"
	"time"
)

// Date represents a day, without time zone consideration
type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*time.Time)(d))
}

const DateLayout = "2006-01-02"

func (st Students) ByIdClassroom() map[int64]Students {
	out := make(map[int64]Students)
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
