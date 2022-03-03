package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	trivialpoursuit "github.com/benoitkugler/maths-online/trivial-poursuit"
)

func main() {
	nbPlayers := flag.Int("players", 2, "number of players per game")
	flag.Parse()

	fmt.Println("Number of players:", *nbPlayers)

	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	u.Path = "trivial-poursuit"

	trivialpoursuit.ProgressLogger.SetOutput(os.Stdout)
	trivialpoursuit.RegisterAndStart("/"+u.Path, trivialpoursuit.GameOptions{PlayersNumber: *nbPlayers})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("incomming")
		rw.Write([]byte("All good"))
	})

	u.Scheme = "ws"
	fmt.Println("Listenning at ", u.String())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
