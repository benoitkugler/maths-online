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

var accessForbidden = errors.New("ressource access forbidden")

// Controller provides the route handling teacher inscription,
// connection and settings.
type Controller struct {
	db                     *sql.DB
	teacherKey, studentKey pass.Encrypter
	smtp                   pass.SMTP
	host                   string // used for links

	admin Teacher // loaded at creation

	classCodes *classroomsCode
}

// NewController return a new controller.
// `LoadAdminTeacher` should be called once.
func NewController(db *sql.DB, smtp pass.SMTP, teacherKey, studentKey pass.Encrypter, host string) *Controller {
	return &Controller{
		db:         db,
		teacherKey: teacherKey,
		studentKey: studentKey,
		smtp:       smtp,
		host:       host,
		classCodes: &classroomsCode{codes: make(map[string]int64)},
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
	_, err := mail.ParseAddress(args.Mail)
	if err != nil {
		return "", errors.New("L'adresse mail est invalide.")
	}

	payload, err := ct.teacherKey.EncryptJSON(args)
	if err != nil {
		return "", err
	}

	url := utils.BuildUrl(ct.host, ValidateInscriptionEndPoint, map[string]string{
		"data": payload,
	})

	return fmt.Sprintf(`
	Bonjour et bienvenue sur Isyro ! <br/><br/>

	Nous avons pu v??rifier la validit?? de votre adresse mail. Merci de terminer votre inscription
	en suivant le lien : <br/>
	<a href="%s">%s</a> <br/><br/>

	Bonne cr??ation p??dagogique ! <br/><br/>

	L'??quipe Isyro
	`, url, url), nil
}

// AskInscription send a link to register a new user account.
func (ct *Controller) AskInscription(c echo.Context) error {
	var args AskInscriptionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.askInscription(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) askInscription(args AskInscriptionIn) (AskInscriptionOut, error) {
	args.Mail = strings.TrimSpace(args.Mail)

	// TODO: should we accept anybody ?

	teachers, err := SelectAllTeachers(ct.db)
	if err != nil {
		return AskInscriptionOut{}, utils.SQLError(err)
	}
	for _, tc := range teachers {
		if tc.Mail == args.Mail {
			return AskInscriptionOut{Error: "Cette adresse mail est d??j?? utilis??e."}, nil
		}
	}

	if len(args.Password) < 2 {
		return AskInscriptionOut{Error: "Merci de choisir un mot de passe plus solide.", IsPasswordError: true}, nil
	}

	mailText, err := ct.emailInscription(args)
	if err != nil {
		return AskInscriptionOut{Error: err.Error()}, nil
	}

	err = mailer.SendMail(ct.smtp, args.Mail, "Bienvenue sur Isyro", mailText)
	if err != nil {
		return AskInscriptionOut{}, fmt.Errorf("Erreur interne (%s)", err)
	}

	return AskInscriptionOut{}, nil
}

func (ct *Controller) ValidateInscription(c echo.Context) error {
	payload := c.QueryParam("data")

	var args AskInscriptionIn
	err := ct.teacherKey.DecryptJSON(payload, &args)
	if err != nil {
		return err
	}

	t := Teacher{
		Mail:            args.Mail,
		PasswordCrypted: ct.teacherKey.EncryptPassword(args.Password),
	}
	t, err = t.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	url := utils.BuildUrl(ct.host, "/prof", map[string]string{
		"show-success-inscription": "OK",
	})
	return c.Redirect(302, url)
}

func (ct *Controller) loggin(args LogginIn) (LogginOut, error) {
	row := ct.db.QueryRow("SELECT * FROM teachers WHERE mail = $1", args.Mail)
	teacher, err := ScanTeacher(row)
	if err == sql.ErrNoRows {
		return LogginOut{Error: "Cette adresse mail n'est pas utilis??e."}, nil
	}
	if err != nil {
		return LogginOut{}, err
	}

	if args.Password != ct.teacherKey.DecryptPassword(teacher.PasswordCrypted) {
		return LogginOut{Error: "Le mot de passe est incorrect.", IsPasswordError: true}, nil
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

type Origin struct {
	AllowPublish bool // is the ressource allowed to be made public ?
	IsPublic     bool // used, for Personnal visibility
	Visibility   Visibility
}
