package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	trivialpoursuit "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate ../../structgen/apigen -source=main.go -out=../prof/src/controller/api_gen.ts

func main() {
	devPtr := flag.Bool("dev", false, "run in dev mode (localhost)")
	dryPtr := flag.Bool("dry", false, "do not listen, but quit early")
	flag.Parse()

	host := getAdress(*devPtr)

	ct := trivialpoursuit.NewController(host)

	e := echo.New()
	e.HideBanner = true

	if *devPtr {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowMethods:  append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions),
			AllowHeaders:  []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
			ExposeHeaders: []string{"Content-Disposition"},
		}))
		fmt.Println("CORS activé.")

		// simulate some latency for better conformity
		e.Use(echo.WrapMiddleware(func(h http.Handler) http.Handler {
			time.Sleep(time.Second / 2)
			return h
		}))
	}

	setupRoutes(e, ct)

	if *dryPtr {
		log.Printf("Setup done, leaving early.")
		return
	}
	fmt.Println("Setup done")

	err := e.Start(host) // start and block
	e.Logger.Fatal(err)  // report error and quit
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

func setupRoutes(e *echo.Echo, ct *trivialpoursuit.Controller) {
	e.POST("/trivial/launch_game", ct.LaunchGame)
	e.GET(trivialpoursuit.GameEndPoint, ct.AccessGame)
}
