package teacher

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestController_mailInscription(t *testing.T) {
	ct := NewController(nil, pass.SMTP{}, pass.Encrypter{}, pass.Encrypter{}, "localhost:1323", "1234")

	if _, err := ct.emailInscription(AskInscriptionIn{Mail: "invalid", Password: "qsqmmùs"}); err == nil {
		t.Fatal("expected error on invalid mail")
	}

	text, err := ct.emailInscription(AskInscriptionIn{
		Mail:     "test@free.fr",
		Password: "my-pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(text)
}

func TestSettingsAPI(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql")
	defer db.Remove()

	ct := Controller{db: db.DB}
	teach := tc.Teacher{
		Mail: "dummy@free.fr",
	}
	teach, err := teach.Insert(ct.db)
	tu.AssertNoErr(t, err)

	settings, err := ct.getSettings(teach.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, settings.Mail == teach.Mail)
	tu.Assert(t, settings.HasEditorSimplified == false)

	err = ct.updateSettings(TeacherSettings{Mail: "anoter@free.fr", HasEditorSimplified: true}, teach.Id)
	tu.AssertNoErr(t, err)
}

func TestPasswordCrypted(t *testing.T) {
	key := pass.Encrypter{}
	password := "hehe"
	fmt.Printf(`'\x%s'`, hex.EncodeToString(key.EncryptPassword(password)))
	fmt.Println()
}
