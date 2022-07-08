package questions

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//go:generate ../../../../../structgen/structgen -source=models.go -mode=itfs-json:gen_itfs.go

type Enonce []Block

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
func (s *QuestionPage) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s QuestionPage) Value() (driver.Value, error) { return dumpJSON(s) }

// Scan implements the driver.Scanner interface using JSON
func (s *Parameters) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s Parameters) Value() (driver.Value, error) { return dumpJSON(s) }
