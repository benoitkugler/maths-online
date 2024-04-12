package teacher

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/mailer"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/reviews"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
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

	admin         tc.Teacher   // loaded at creation
	demoClassroom tc.Classroom // loaded at creation

	demoCode string
}

// NewController return a new controller.
// `LoadAdminTeacher` should be called once.
func NewController(db *sql.DB, smtp pass.SMTP, teacherKey, studentKey pass.Encrypter, host, demoCode string) *Controller {
	return &Controller{
		db:         db,
		teacherKey: teacherKey,
		studentKey: studentKey,
		smtp:       smtp,
		host:       host,
		demoCode:   demoCode,
	}
}

// LoadAdminTeacher loads and stores the admin account.
// By convention, only one account has admin rights. It is manually created at
// DB setup, and never added (neiter removed) at run time.
func (ct *Controller) LoadAdminTeacher() (tc.Teacher, error) {
	rows, err := ct.db.Query("SELECT * FROM teachers WHERE IsAdmin = true")
	if err != nil {
		return tc.Teacher{}, utils.SQLError(err)
	}
	teachers, err := tc.ScanTeachers(rows)
	if err != nil {
		return tc.Teacher{}, utils.SQLError(err)
	}
	if len(teachers) != 1 {
		return tc.Teacher{}, errors.New("internal error: exactly one teacher must be admin")
	}
	ct.admin = teachers[teachers.IDs()[0]]

	return ct.admin, nil
}

// LoadDemoClassroom loads and stores the demo classroom, which is a [Classroom]
// manually created at DB setup, and attributed to the admin account,
// with a special ID = 1 .
func (ct *Controller) LoadDemoClassroom() (tc.Classroom, error) {
	cl, err := tc.SelectClassroom(ct.db, 1)
	if err != nil {
		return tc.Classroom{}, utils.SQLError(err)
	}
	// sanity checks
	if cl.IdTeacher != ct.admin.Id {
		return tc.Classroom{}, errors.New("internal error: unexpected owner of the the demo classroom")
	}
	ct.demoClassroom = cl

	return cl, nil
}

func (ct *Controller) CleanupClassroomCodes() error { return tc.CleanupClassroomCodes(ct.db) }

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

	Nous avons pu vérifier la validité de votre adresse mail. Merci de terminer votre inscription
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

	out, err := ct.askInscription(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) askInscription(args AskInscriptionIn) (AskInscriptionOut, error) {
	args.Mail = strings.TrimSpace(args.Mail)

	// TODO: should we accept anybody ?

	teachers, err := tc.SelectAllTeachers(ct.db)
	if err != nil {
		return AskInscriptionOut{}, utils.SQLError(err)
	}
	for _, tc := range teachers {
		if tc.Mail == args.Mail {
			return AskInscriptionOut{Error: "Cette adresse mail est déjà utilisée."}, nil
		}
	}

	if len(args.Password) < 2 {
		return AskInscriptionOut{Error: "Merci de choisir un mot de passe plus solide.", IsPasswordError: true}, nil
	}

	mailText, err := ct.emailInscription(args)
	if err != nil {
		return AskInscriptionOut{Error: err.Error()}, nil
	}

	err = mailer.SendMail(ct.smtp, []string{args.Mail}, "Bienvenue sur Isyro", mailText)
	if err != nil {
		return AskInscriptionOut{}, fmt.Errorf("Erreur interne (%s)", err)
	}

	return AskInscriptionOut{}, nil
}

func hasEditorSimplified(topic tc.MatiereTag) bool {
	switch topic {
	case tc.Mathematiques, tc.PhysiqueChimie, tc.SVT, tc.SES:
		return false
	default:
		return true
	}
}

func (ct *Controller) ValidateInscription(c echo.Context) error {
	payload := c.QueryParam("data")

	var args AskInscriptionIn
	err := ct.teacherKey.DecryptJSON(payload, &args)
	if err != nil {
		return err
	}

	t := tc.Teacher{
		Mail:                args.Mail,
		PasswordCrypted:     ct.teacherKey.EncryptPassword(args.Password),
		FavoriteMatiere:     args.FavoriteMatiere,
		HasSimplifiedEditor: hasEditorSimplified(args.FavoriteMatiere),
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
	teacher, err := tc.ScanTeacher(row)
	if err == sql.ErrNoRows {
		return LogginOut{Error: "Cette adresse mail n'est pas utilisée."}, nil
	}
	if err != nil {
		return LogginOut{}, err
	}

	if args.Password != ct.teacherKey.DecryptPassword(teacher.PasswordCrypted) {
		return LogginOut{Error: "Le mot de passe est incorrect.", IsPasswordError: true}, nil
	}

	token, err := ct.newToken(teacher.Id)
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

// TeacherResetPassword generates a new password for the given account
// and sends it by email.
func (ct *Controller) TeacherResetPassword(c echo.Context) error {
	mail := c.QueryParam("mail")
	err := ct.resetPassword(mail)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func (ct *Controller) emailResetPassword(newPassword string) string {
	url := utils.BuildUrl(ct.host, "prof", nil)
	return fmt.Sprintf(`
	Bonjour, <br/><br/>

	Vous avez demandé la ré-initialisation de votre mot de passe Isyro. Votre nouveau mot de passe est : <br/>
	<b>%s</b> <br/><br/>

	Après vous être <a href="%s">connecté</a>, vous pourrez le modifier dans vos réglages.<br/><br/>

	Bonne création pédagogique ! <br/><br/>

	L'équipe Isyro`, newPassword, url)
}

func (ct *Controller) resetPassword(mail string) error {
	row := ct.db.QueryRow("SELECT * FROM teachers WHERE mail = $1", mail)
	teacher, err := tc.ScanTeacher(row)
	if err == sql.ErrNoRows {
		return errors.New("Cette adresse mail n'est pas utilisée.")
	}
	if err != nil {
		return utils.SQLError(err)
	}
	// generate a new password
	newPassword := utils.RandomString(true, 8)
	teacher.PasswordCrypted = ct.teacherKey.EncryptPassword(newPassword)
	_, err = teacher.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// send it by email
	mailText := ct.emailResetPassword(newPassword)
	err = mailer.SendMail(ct.smtp, []string{mail}, "Changement de mot de passe", mailText)
	if err != nil {
		return fmt.Errorf("Erreur interne (%s)", err)
	}

	return nil
}

type TeacherSettings struct {
	Mail                string
	Password            string
	HasEditorSimplified bool
	Contact             tc.Contact
	FavoriteMatiere     teacher.MatiereTag
}

// TeacherGetSettings returns the teacher global settings.
func (ct *Controller) TeacherGetSettings(c echo.Context) error {
	userID := JWTTeacher(c)

	out, err := ct.getSettings(userID)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

func (ct *Controller) getSettings(userID teacher.IdTeacher) (TeacherSettings, error) {
	teach, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return TeacherSettings{}, utils.SQLError(err)
	}

	password := ct.teacherKey.DecryptPassword(teach.PasswordCrypted)

	return TeacherSettings{
		Mail:                teach.Mail,
		Password:            password,
		HasEditorSimplified: teach.HasSimplifiedEditor,
		Contact:             teach.Contact,
		FavoriteMatiere:     teach.FavoriteMatiere,
	}, nil
}

func (ct *Controller) TeacherUpdateSettings(c echo.Context) error {
	userID := JWTTeacher(c)

	var args TeacherSettings
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateSettings(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSettings(args TeacherSettings, userID teacher.IdTeacher) error {
	teach, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return utils.SQLError(err)
	}

	if len(args.Password) < 3 {
		return errors.New("internal error: password too short")
	}

	teach.Mail = args.Mail
	teach.PasswordCrypted = ct.teacherKey.EncryptPassword(args.Password)
	teach.HasSimplifiedEditor = args.HasEditorSimplified
	teach.Contact = args.Contact
	teach.FavoriteMatiere = args.FavoriteMatiere

	_, err = teach.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// shared types

type Origin struct {
	Visibility   Visibility
	PublicStatus PublicStatus
	IsInReview   OptionalIdReview // true if the owner has already started a review for the resource
}

type OptionalIdReview struct {
	InReview bool
	Id       reviews.IdReview
}
