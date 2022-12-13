package teacher

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/pass"
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

func TestPasswordCrypted(t *testing.T) {
	key := pass.Encrypter{}
	password := "hehe"
	fmt.Printf(`'\x%s'`, hex.EncodeToString(key.EncryptPassword(password)))
	fmt.Println()
}
