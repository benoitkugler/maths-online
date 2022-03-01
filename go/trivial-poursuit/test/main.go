package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	trivialpoursuit "github.com/benoitkugler/maths-online/trivial-poursuit"
)

func main() {
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	u.Path = "trivial-poursuit"
	trivialpoursuit.RegisterAndStart("/"+u.Path, trivialpoursuit.GameOptions{PlayersNumber: 1})

	u.Scheme = "ws"
	fmt.Println("Listenning at ", u.String())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
