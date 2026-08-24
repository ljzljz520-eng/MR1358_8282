package main

import (
	"example.com/nightguide/internal/app"
	"flag"
	"log"
	"net/http"
)

func main() {
	path := flag.String("db", "nightguide.db", "sqlite database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	application, err := app.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()
	log.Printf("nightguide listening on %s", *addr)
	if err := http.ListenAndServe(*addr, application.HTTP.Handler()); err != nil {
		log.Fatal(err)
	}
}
