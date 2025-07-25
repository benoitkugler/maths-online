package teacher

import (
	"time"

	"github.com/benoitkugler/maths-online/server/src/utils"
)

const DateLayout = "2006-01-02"

type Time time.Time

func (d Time) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Time) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

// Date represents a day, without time zone consideration
type Date time.Time

func (d Date) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Date) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func NewDateFrom(ti time.Time) Date {
	return NewDate(ti.Year(), ti.Month(), ti.Day())
}

func (d Date) Time() time.Time {
	ti := time.Time(d)
	return time.Date(ti.Year(), ti.Month(), ti.Day(), 0, 0, 0, 0, time.UTC)
}

type Contact struct {
	Name string
	URL  string
}

// MatiereTag are special tags to indicate the topic of a question/exercice
type MatiereTag string

const (
	Autre          MatiereTag = "AUTRE"
	Mathematiques  MatiereTag = "MATHS"
	Francais       MatiereTag = "FRANCAIS"
	HistoireGeo    MatiereTag = "HISTOIRE-GEO"
	Espagnol       MatiereTag = "ESPAGNOL"
	Italien        MatiereTag = "ITALIEN"
	Allemand       MatiereTag = "ALLEMAND"
	Anglais        MatiereTag = "ANGLAIS"
	SES            MatiereTag = "SES"
	PhysiqueChimie MatiereTag = "PHYSIQUE"
	SVT            MatiereTag = "SVT"
)

// CleanupClassroomCodes removes the expired codes.
func CleanupClassroomCodes(db DB) error {
	_, err := db.Exec("DELETE FROM classroom_codes WHERE now() > expiresat;")
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

func (ccs ClassroomCodes) Codes() map[string]bool {
	out := make(map[string]bool)
	for _, item := range ccs {
		out[item.Code] = true
	}
	return out
}

type Client struct {
	Device string    // the name of the device
	Time   time.Time // when the client was connected
}

type Clients []Client
