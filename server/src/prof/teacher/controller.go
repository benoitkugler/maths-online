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

	admin Teacher // loaded at creation
}

// NewController return a new controller.
// `LoadAdminTeacher` should be called once.
func NewController(db *sql.DB, smtp pass.SMTP, key pass.Encrypter, host string) *Controller {
	return &Controller{
		db:   db,
		key:  key,
		smtp: smtp,
		host: host,
	}
}

// LoadAdminTeacher loads and stores the admin account.
// By convention, only one account has admin rights. It is manually created at
// DB setup, and never added (neiter removed) at run time.
func (ct *Controller) LoadAdminTeacher() (Teacher, error) {
	rows, err := ct.db.Query("SELECT * FROM teachers WHERE is_admin = true")
	if err != nil {
		return Teacher{}, utils.SQLError(err)
	}
	teachers, err := ScanTeachers(rows)
	if err != nil {
		return Teacher{}, utils.SQLError(err)
	}
	if len(teachers) != 1 {
		return Teacher{}, errors.New("exactly one teacher must be admin")
	}
	ct.admin = teachers[teachers.IDs()[0]]
	return ct.admin, nil
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

func (ct *Controller) loggin(args LogginIn) (LogginOut, error) {
	row := ct.db.QueryRow("SELECT * FROM teachers WHERE mail = $1", args.Mail)
	teacher, err := ScanTeacher(row)
	if err == sql.ErrNoRows {
		return LogginOut{Error: "Cette adresse mail n'est pas utilisée."}, nil
	}
	if err != nil {
		return LogginOut{}, err
	}

	if args.Password != ct.key.DecryptPassword(teacher.PasswordCrypted) {
		return LogginOut{Error: "Le mot de passe est incorrect."}, nil
	}

	token, err := ct.newToken(teacher)
	if err != nil {
		return LogginOut{}, err
	}

	return LogginOut{Token: token}, nil
}

func (ct *Controller) Loggin(c echo.Context) error {
	var args LogginIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.loggin(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}
