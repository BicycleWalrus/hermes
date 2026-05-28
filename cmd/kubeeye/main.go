package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BicycleWalrus/hermes/pkg/kubeeye"
)

const defaultAddr = ":8888"

func main() {
	addr := os.Getenv("KUBEEYE_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           kubeeye.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("kubeeye listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
