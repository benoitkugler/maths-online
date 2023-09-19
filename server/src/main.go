package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/mailer"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/prof/editor"
	"github.com/benoitkugler/maths-online/server/src/prof/homework"
	"github.com/benoitkugler/maths-online/server/src/prof/reviews"
	"github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/prof/trivial"
	tvGame "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/vitrine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate ../../../gomacro/cmd/gomacro gomacro.json

func connectDB(dev bool) (*sql.DB, error) {
	var credentials pass.DB
	if dev {
		credentials = pass.DB{
			Host:     "localhost",
			User:     "benoit",
			Password: "dummy",
			Name:     "isyro_prod",
		}
	} else { // in production, read from env
		var err error
		credentials, err = pass.NewDB()
		if err != nil {
			return nil, err
		}
	}

	db, err := credentials.ConnectPostgres()
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Printf("DB configured with %v connected.\n", credentials)

	return db, err
}

func getStudentEncrypter(dev bool) (out pass.Encrypter) {
	if dev {
		return pass.Encrypter{1, 2, 3, 4, 5, 6}
	} else {
		out, err := pass.NewEncrypter("STUDENT_ENC_KEY")
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
}

func getTeacherEncrypter(dev bool) (out pass.Encrypter) {
	if dev {
		return pass.Encrypter{4, 5, 6, 7, 8, 9}
	} else {
		out, err := pass.NewEncrypter("TEACHER_ENC_KEY")
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
}

func getDemoCode() string {
	demoCode := os.Getenv("DEMO_CODE")
	if demoCode == "" {
		log.Fatal("Missing DEMO_CODE env. variable")
	}
	return demoCode
}

func getAdminEmails() []string {
	m := os.Getenv("ADMIN_MAILS")
	if m == "" {
		log.Fatal("Missing ADMIN_MAILS env. variable")
	}
	return strings.Split(os.Getenv("ADMIN_MAILS"), ",")
}

func devSetup(e *echo.Echo, tc *teacher.Controller) {
	dev, err := tc.GetDevToken()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dev)

	// also Cross origin requests
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:  append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions),
		AllowHeaders:  []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		ExposeHeaders: []string{"Content-Disposition"},
	}))
	fmt.Println("CORS activé.")

	devMail := os.Getenv("DEV_MAIL_TO")
	if devMail == "" {
		log.Fatal("Missing env. variable DEV_MAIL_TO")
	}
	mailer.SetDevMail(devMail)
	fmt.Println("Mail redirected to ", devMail)
}

func sanityChecks(db *sql.DB, skipValidation bool) {
	if skipValidation {
		fmt.Println("Validation skipped.")
		return
	}
	if err := editor.ValidateAllQuestions(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Question table validated.")
}

func main() {
	devPtr := flag.Bool("dev", false, "run in dev mode (localhost)")
	dryPtr := flag.Bool("dry", false, "do not listen, but quit early")
	skipValidation := flag.Bool("s", false, "skip question table validation")
	flag.Parse()

	adress := getAdress(*devPtr)
	host := getPublicHost(*devPtr)

	studentKey := getStudentEncrypter(*devPtr)
	fmt.Printf("Student encrypter setup with key: %v\n", studentKey)

	teacherKey := getTeacherEncrypter(*devPtr)
	fmt.Printf("Teacher encrypter setup with key: %v\n", teacherKey)

	demoCode := getDemoCode()
	fmt.Printf("Demontration code for student activities: %s\n", demoCode)

	adminEmails := getAdminEmails()
	fmt.Printf("Admin emails for contact form: %v\n", adminEmails)

	db, err := connectDB(*devPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	smtp, err := pass.NewSMTP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SMTP configured with %v.\n", smtp)

	tc := teacher.NewController(db, smtp, teacherKey, studentKey, host, demoCode)
	admin, err := tc.LoadAdminTeacher()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Admin teacher loaded.")
	_, err = tc.LoadDemoClassroom()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Demo classroom loaded.")
	err = tc.CleanupClassroomCodes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Expired classroom codes cleaned up.")

	tvc := trivial.NewController(db, studentKey, demoCode, admin)
	if err = tvc.CheckDemoQuestions(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Trivial demo questions checked.")

	hwc := homework.NewController(db, admin, studentKey)
	edit := editor.NewController(db, admin)
	vit := &vitrine.Controller{Smtp: smtp, AdminMails: adminEmails}
	review := reviews.NewController(db, admin, smtp)

	// for now, show the logs
	tvGame.ProgressLogger.SetOutput(os.Stdout)
	tvGame.WarningLogger.SetOutput(os.Stdout)

	e := echo.New()
	e.HideBanner = true
	// this prints detailed error messages
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		err = echo.NewHTTPError(400, err.Error())
		e.DefaultHTTPErrorHandler(err, c)
	}

	if *devPtr {
		devSetup(e, tc)
	}

	setupRoutes(e, db, tvc, edit, tc, hwc, vit, review)

	if *dryPtr {
		sanityChecks(db, *skipValidation)
		fmt.Println("Setup done, leaving early.")
		return
	} else {
		go sanityChecks(db, *skipValidation)
	}
	fmt.Println("Setup done (pending sanityChecks)")

	err = e.Start(adress) // start and block
	e.Logger.Fatal(err)   // report error and quit
}

func getPublicHost(dev bool) string {
	if dev {
		return "localhost:1323"
	}
	// use env variable
	host := os.Getenv("PUBLIC_HOST")
	if host == "" {
		log.Fatal("misssing PUBLIC_HOST env variable")
	}
	return host
}

func getAdress(dev bool) string {
	var adress string
	if dev {
		adress = "localhost:1323"
	} else {
		// alwaysdata use IP and PORT env var
		host := os.Getenv("IP")
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatal("No PORT found ", err)
		}
		if strings.Count(host, ":") >= 2 { // ipV6 -> besoin de crochet
			host = "[" + host + "]"
		}
		adress = fmt.Sprintf("%s:%d", host, port)
	}
	return adress
}

// noCache prevent the browser to cache the file served,
// so that the build frontend app are always up to date.
func noCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store")
		c.Response().Header().Set("Expires", "0")
		return next(c)
	}
}

// cacheIframe set a short cache time
func cacheIframe(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=public,21600")
		return next(c)
	}
}

// cacheStatic adopt a very aggressive caching policy, suitable
// for immutable content
func cacheStatic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=31536000")
		return next(c)
	}
}

func serveVitrineApp(c echo.Context) error {
	return c.File("static/vitrine/index.html")
}

func serveProfApp(c echo.Context) error {
	return c.File("static/prof/index.html")
}

func serveProfLoopbackApp(c echo.Context) error {
	return c.File("static/prof_loopback/index.html")
}

func serveProfPreviewApp(c echo.Context) error {
	return c.File("static/prof_preview/index.html")
}

func serveEleveApp(c echo.Context) error {
	return c.File("static/eleve/index.html")
}

func setupRoutes(e *echo.Echo, db *sql.DB,
	tvc *trivial.Controller, edit *editor.Controller,
	tc *teacher.Controller, home *homework.Controller,
	vit *vitrine.Controller, review *reviews.Controller,
) {
	setupProfAPI(e, tvc, edit, tc, home, review)

	// main page
	e.GET("", serveVitrineApp)
	e.Static("/*", "static/vitrine")
	e.POST("/reponse-contact", vit.HandleFormMessage)

	// global static files used by frontend apps
	e.Group("/static", middleware.Gzip(), cacheStatic).Static("/*", "static")

	e.GET("/test-eleve", serveEleveApp, noCache)
	e.GET("/test-eleve/", serveEleveApp, noCache)
	e.Group("/test-eleve/*", middleware.Gzip(), cacheStatic).Static("/*", "static/eleve")

	e.GET("/prof-loopback-app", serveProfLoopbackApp, cacheIframe)
	e.GET("/prof-loopback-app/", serveProfLoopbackApp, cacheIframe)
	e.Group("/prof-loopback-app/*", middleware.Gzip(), cacheStatic).Static("/*", "static/prof_loopback")

	e.GET("/prof-preview-app", serveProfPreviewApp, cacheIframe)
	e.GET("/prof-preview-app/", serveProfPreviewApp, cacheIframe)
	e.Group("/prof-preview-app/*", middleware.Gzip(), cacheStatic).Static("/*", "static/prof_preview")

	// student trivial access
	e.GET("/trivial/game/setup", tvc.SetupStudentClient)
	e.GET("/trivial/game/connect", tvc.ConnectStudentSession)
	// student trivial self access launcher
	e.GET("/api/student/trivial/selfaccess", tvc.StudentGetSelfaccess)
	e.GET("/api/student/trivial/selfaccess/launch", tvc.StudentLaunchSelfaccess)
	e.GET("/api/student/trivial/selfaccess/start", tvc.StudentStartSelfaccess)
	// trivial monitor
	e.GET("/api/trivial/monitor", tvc.GetTrivialsMetrics)

	// student client classroom managment
	e.GET("/api/classroom/login", tc.CheckStudentClassroom)
	e.GET("/api/classroom/attach", tc.AttachStudentToClassroomStep1)
	e.POST("/api/classroom/attach", tc.AttachStudentToClassroomStep2)

	// prof. back office
	for _, route := range []string{
		"/prof",
		"/prof/*",
	} {
		e.GET(route, serveProfApp, noCache)
	}

	// embeded preview app
	e.POST("/api/loopack/evaluate-question", edit.LoopackEvaluateQuestion)
	e.POST("/api/loopack/question-answer", edit.LoopbackShowQuestionAnswer)

	// shared expression syntax check endpoint
	e.GET("/api/check-expression", checkExpressionSyntax)
	e.POST("/api/evaluate-question", func(c echo.Context) error {
		return evaluateQuestion(db, c)
	})

	// standalone question/exercice
	e.POST("/api/questions/instantiate", func(c echo.Context) error {
		return instantiateQuestions(db, c)
	})
	e.POST("/api/questions/evaluate", func(c echo.Context) error {
		return evaluateQuestion(db, c)
	})
	e.POST("/api/exercices/evaluate", func(c echo.Context) error {
		return evaluateExercice(db, c)
	})

	// student homework API
	e.GET("/api/student/homework/sheets", home.StudentGetTravaux)
	e.GET("/api/student/homework/sheets/free", home.StudentGetFreeTravaux)
	e.GET("/api/student/homework/task/instantiate", home.StudentInstantiateTask)
	e.POST("/api/student/homework/task/evaluate", home.StudentEvaluateTask)
	e.POST("/api/student/homework/task/reset", home.StudentResetTask)
}
