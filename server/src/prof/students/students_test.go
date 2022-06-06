package students

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestTime(t *testing.T) {
	ti := time.Now()
	if ti.String()[0:10] != ti.Format(DateLayout) {
		t.Fatal()
	}
}

func TestSQLTime(t *testing.T) {
	db := testutils.CreateDBDev(t, "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	st, err := Student{Birthday: Date(time.Now())}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	if time.Time(st.Birthday).IsZero() {
		t.Fatal()
	}
}
