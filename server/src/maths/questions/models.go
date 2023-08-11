package questions

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
)

// this file is used to generate JSON routines for interfaces

// Parameters stores the definition of the random parameters, in the order
// used by the user
type Parameters []ParameterEntry

type Enonce []Block

type GeometricConstructionFieldBlock struct {
	Field      GeoField
	Background FiguresOrGraphs
}

type GeoField interface {
	instantiate(params expression.Vars) (geoFieldInstance, error)
	setupValidator(params expression.RandomParameters) (validator, error)
}
type FiguresOrGraphs interface {
	instantiateFG(params expression.Vars) (client.FigureOrGraph, error)
	setupValidator(params expression.RandomParameters) (validator, error)
}

func loadJSON(out interface{}, src interface{}) error {
	if src == nil {
		return nil // zero value out
	}
	bs, ok := src.([]byte)
	if !ok {
		return errors.New("not a []byte")
	}
	return json.Unmarshal(bs, out)
}

func dumpJSON(s interface{}) (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return driver.Value(string(b)), nil
}

// Scan implements the driver.Scanner interface using JSON
func (s *Enonce) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s Enonce) Value() (driver.Value, error) { return dumpJSON(s) }

// Scan implements the driver.Scanner interface using JSON
func (s *Parameters) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s Parameters) Value() (driver.Value, error) { return dumpJSON(s) }
