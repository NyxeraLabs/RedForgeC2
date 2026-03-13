package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"redforgec2/server/go_backend/internal/httpapi"
	"redforgec2/server/go_backend/internal/store"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address (lab-only; default binds to localhost)")
	flag.Parse()

	st := store.New()
	router := httpapi.NewRouter(st)

	srv := &http.Server{
		Addr:              *addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("redforgec2 teamserver (simulation) listening on http://%s", *addr)
	log.Fatal(srv.ListenAndServe())
}

