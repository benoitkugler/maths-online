package editor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/pass"
)

func TestInstantiateQuestions(t *testing.T) {
	creds := pass.DB{
		Host:     "localhost",
		User:     "benoit",
		Password: "dummy",
		Name:     "isyro_prod",
	}
	db, err := creds.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", creds, err)
		return
	}

	ct := NewController(db)
	out, err := ct.InstantiateQuestions([]int64{24, 29, 37})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := json.MarshalIndent(out, " ", " ")
	fmt.Println(string(s))
}
