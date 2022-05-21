package teacher

import (
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/pass"
)

func TestController_mailInscription(t *testing.T) {
	ct := NewController(nil, pass.SMTP{}, pass.Encrypter{}, "localhost:1323")

	if _, err := ct.emailInscription(AskInscriptionIn{Mail: "invalid", Password: "qsqmmùs"}); err == nil {
		t.Fatal("expected error on invalid mail")
	}

	if _, err := ct.emailInscription(AskInscriptionIn{Mail: "ok@free.fr", Password: ""}); err == nil {
		t.Fatal("expected error on empty password")
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
