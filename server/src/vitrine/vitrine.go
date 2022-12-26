// Pacakge vitrine provides the logic required
// to handle interactivity on the landing page
package vitrine

import (
	"fmt"
	"log"
	"time"

	"github.com/benoitkugler/maths-online/server/src/mailer"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	Smtp       pass.SMTP
	AdminMails []string
}

func (ct *Controller) HandleFormMessage(c echo.Context) error {
	var form struct {
		Mail    string `form:"mail"`
		Message string `form:"message"`
	}
	if err := c.Bind(&form); err != nil {
		return err
	}

	const timeLayout = "le 02/01/2006 à 15h04"

	body := fmt.Sprintf(`Message reçu depuis le formulaire de contact, %s.<br/>
	
	Message : <br/>
	%s
	
	<br/>
	<br/>
	Adresse mail : %s
	`, time.Now().Format(timeLayout), form.Message, form.Mail)

	log.Println(body)

	if len(ct.AdminMails) != 0 {
		err := mailer.SendMail(ct.Smtp, ct.AdminMails, "[Isyro] - Nouveau message", body)
		if err != nil {
			return err
		}
	}

	return c.File("static/vitrine/reponse-contact.html")
}
