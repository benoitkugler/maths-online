package teacher

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/benoitkugler/maths-online/mailer"
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

// Controller provides the route handling teacher inscription,
// connection and settings.
type Controller struct {
	db   *sql.DB
	key  pass.Encrypter
	smtp pass.SMTP
	host string // used for links
}

func NewController(db *sql.DB, smtp pass.SMTP, key pass.Encrypter, host string) *Controller {
	return &Controller{
		db:   db,
		key:  key,
		smtp: smtp,
		host: host,
	}
}

const ValidateInscriptionEndPoint = "inscription"

func (ct *Controller) emailInscription(args AskInscriptionIn) (string, error) {
	if len(args.Password) < 2 {
		return "", errors.New("Merci de choisir un mot de passe plus solide.")
	}

	_, err := mail.ParseAddress(args.Mail)
	if err != nil {
		return "", errors.New("L'adresse mail est invalide.")
	}

	payload, err := ct.key.EncryptJSON(args)
	if err != nil {
		return "", err
	}

	url := utils.BuildUrl(ct.host, ValidateInscriptionEndPoint, map[string]string{
		"data": payload,
	})

	return fmt.Sprintf(`
	Bonjour et bienvenue sur Isyro ! <br/><br/>

	Votre adresse mail est bien valide. Merci de terminer votre inscription
	en suivant le lien : <br/>
	<a href="%s">%s</a> <br/><br/>

	Bonne création pédagogique ! <br/><br/>

	L'équipe Isyro
	`, url, url), nil
}

// AskInscription send a link to register a new user account.
func (ct *Controller) AskInscription(c echo.Context) error {
	var args AskInscriptionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	args.Mail = strings.TrimSpace(args.Mail)

	// TODO: should we accept anybody ?

	teachers, err := SelectAllTeachers(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}
	for _, tc := range teachers {
		if tc.Mail == args.Mail {
			return errors.New("Cette adresse mail est déjà utilisée.")
		}
	}

	mailText, err := ct.emailInscription(args)
	if err != nil {
		return err
	}

	err = mailer.SendMail(ct.smtp, args.Mail, "Bienvenue sur Isyro", mailText)
	if err != nil {
		return fmt.Errorf("Erreur interne (%s)", err)
	}

	return c.NoContent(200)
}

func (ct *Controller) ValidateInscription(c echo.Context) error {
	payload := c.QueryParam("data")

	var args AskInscriptionIn
	err := ct.key.DecryptJSON(payload, &args)
	if err != nil {
		return err
	}

	t := Teacher{
		Mail:            args.Mail,
		PasswordCrypted: ct.key.EncryptPassword(args.Password),
	}
	t, err = t.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// TODO:
	return c.HTML(200, "OK")
}
