package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	email := os.Getenv("CF_API_EMAIL")
	key := os.Getenv("CF_API_KEY")

	if email == "" || key == "" {
		log.Fatal("CF_API_EMAIL or CF_API_KEY missing from environment")
	}

	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", &updater{email, key}); err != nil {
		log.Println(err)
	}
}

type updater struct {
	email, key string
}

func (u *updater) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	log.Println("Update DNS ", q.Get("hostname"), " to ", q.Get("myip"))
	w.WriteHeader(500)
}
