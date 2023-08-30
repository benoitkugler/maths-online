package teacher

import "time"

const DateLayout = "2006-01-02"

// Date represents a day, without time zone consideration
type Date time.Time

func (d Date) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Date) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

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
