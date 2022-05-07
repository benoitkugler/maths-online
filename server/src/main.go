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

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/editor"
	trivialpoursuit "github.com/benoitkugler/maths-online/prof/trivial-poursuit"
	tvGame "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func connectDB(dev bool) (*sql.DB, error) {
	var credentials pass.DB
	if dev {
		credentials = pass.DB{
			Host:     "localhost",
			User:     "benoit",
			Password: "dummy",
			// Name:     "maths_dev",
			Name: "isyro_prod",
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

func getEncrypter(dev bool) (out pass.Encrypter, err error) {
	if dev {
		out = pass.Encrypter{1, 2, 3, 4, 5, 6}
	} else {
		out, err = pass.NewEncrypter("ENC_KEY")
	}

	fmt.Printf("Encrypter setup with key %v.\n", out)
	return out, err
}

func main() {
	devPtr := flag.Bool("dev", false, "run in dev mode (localhost)")
	dryPtr := flag.Bool("dry", false, "do not listen, but quit early")
	flag.Parse()

	host := getAdress(*devPtr)

	key, err := getEncrypter(*devPtr)
	if err != nil {
		log.Fatal(err)
	}

	db, err := connectDB(*devPtr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go func() {
		if err = exercice.ValidateAllQuestions(db); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Question table validated.")
	}()

	trivial := trivialpoursuit.NewController(db, key)

	// for now, show the logs
	tvGame.ProgressLogger.SetOutput(os.Stdout)
	tvGame.WarningLogger.SetOutput(os.Stdout)

	edit := editor.NewController(db)

	e := echo.New()
	e.HideBanner = true
	// this prints detailed error messages
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		err = echo.NewHTTPError(400, err.Error())
		e.DefaultHTTPErrorHandler(err, c)
	}

	if *devPtr {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowMethods:  append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions),
			AllowHeaders:  []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
			ExposeHeaders: []string{"Content-Disposition"},
		}))
		fmt.Println("CORS activé.")
	}

	setupRoutes(e, trivial, edit)

	if *dryPtr {
		fmt.Println("Setup done, leaving early.")
		return
	}
	fmt.Println("Setup done")

	err = e.Start(host) // start and block
	e.Logger.Fatal(err) // report error and quit
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

// cacheStatic adopt a very aggressive caching policy, suitable
// for immutable content
func cacheStatic(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=31536000")
		return next(c)
	}
}

func serveProfApp(c echo.Context) error {
	return c.File("static/prof/index.html")
}

func serveProfLoopbackApp(c echo.Context) error {
	return c.File("static/prof_loopback/index.html")
}

func serveEleveApp(c echo.Context) error {
	return c.File("static/eleve/index.html")
}

func setupRoutes(e *echo.Echo, trivial *trivialpoursuit.Controller, edit *editor.Controller) {
	setupProfAPI(e, trivial, edit)
	// to sync with the client navigator.sendBeacon
	e.POST("/prof/editor/api/end-preview/:sessionID", edit.EditorEndPreview)

	// global static files used by frontend apps
	e.Group("/static", middleware.Gzip(), cacheStatic).Static("/*", "static")

	e.GET("/test-eleve", serveEleveApp, noCache)
	e.GET("/test-eleve/", serveEleveApp, noCache)
	e.Group("/test-eleve/*", middleware.Gzip(), cacheStatic).Static("/*", "static/eleve")

	e.GET("/prof-loopback-app", serveProfLoopbackApp, noCache)
	e.GET("/prof-loopback-app/", serveProfLoopbackApp, noCache)
	e.Group("/prof-loopback-app/*", middleware.Gzip(), cacheStatic).Static("/*", "static/prof_loopback")

	// e.GET("/trivial/stats", trivial.ShowStats)
	e.GET("/prof/trivial/monitor", trivial.ConnectTeacherMonitor)
	e.GET(trivialpoursuit.GameEndPoint, trivial.ConnectStudentSession)

	// prof. back office
	for _, route := range []string{
		"/prof",
		"/prof/",
		"/prof/trivial",
		"/prof/trivial/",
		"/prof/editor",
		"/prof/editor/",
	} {
		e.GET(route, serveProfApp, noCache)
	}

	// embeded preview app
	e.GET(editor.LoopbackEndpoint, edit.AccessLoopback)

	// shared expression syntax check endpoint
	e.GET("/api/check-expression", checkExpressionSyntax)
	e.POST("/api/evaluate-question", func(c echo.Context) error {
		return evaluateQuestion(edit, c)
	})

	// temporary question quick access

	e.GET("/questions", func(c echo.Context) error {
		var out []client.Question
		for _, qu := range exercice.PredefinedQuestions {
			out = append(out, qu.ToClient())
		}
		return c.JSON(200, out)
	})

	e.POST("/questions/syntaxe/:index", func(c echo.Context) error {
		index, _ := strconv.Atoi(c.Param("index"))
		var data client.QuestionSyntaxCheckIn
		err := c.Bind(&data)
		if err != nil {
			return err
		}

		out := exercice.PredefinedQuestions[index].CheckSyntaxe(data)

		return c.JSON(200, out)
	})

	e.POST("/questions/answer/:index", func(c echo.Context) error {
		index, _ := strconv.Atoi(c.Param("index"))

		var data client.QuestionAnswersIn
		err := c.Bind(&data)
		if err != nil {
			return err
		}

		out := exercice.PredefinedQuestions[index].EvaluateAnswer(data)

		return c.JSON(200, out)
	})
}
