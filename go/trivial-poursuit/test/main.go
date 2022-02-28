package main

import (
	"fmt"
	"log"
	"net/http"

	trivialpoursuit "github.com/benoitkugler/maths-online/trivial-poursuit"
)

func main() {
	trivialpoursuit.SetupRoutes("/ws")

	fmt.Println("Listenning...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
