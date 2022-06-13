package students

import (
	"encoding/json"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*time.Time)(d))
}

const DateLayout = "2006-01-02"
