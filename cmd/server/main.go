package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"asky/internal/router"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "API server port")
	flag.Parse()

	r := router.New()

	log.Printf("Click: http://127.0.0.1:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}