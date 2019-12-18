package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(dump)); err != nil {
		log.Println(err)
	}
}

func dump(_ http.ResponseWriter, r *http.Request) {
	log.Println("From: ", r.RemoteAddr)
	log.Println(r.Method, r.Host, r.RequestURI)
	log.Println(r.BasicAuth())
	log.Println(r.UserAgent())
	log.Println(r.Header)
	_, _ = io.Copy(os.Stdout, r.Body)
	log.Println()
}
