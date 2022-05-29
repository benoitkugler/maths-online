package editor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestInstantiateQuestions(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, teacher.Teacher{})
	out, err := ct.InstantiateQuestions([]int64{24, 29, 37})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := json.MarshalIndent(out, " ", " ")
	fmt.Println(string(s))
}
