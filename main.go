package main

import (
	"./cf"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	email := os.Getenv("CF_API_EMAIL")
	key := os.Getenv("CF_API_KEY")

	if email == "" || key == "" {
		log.Fatal("CF_API_EMAIL or CF_API_KEY missing from environment")
	}

	client, err := cf.NewClient(email, key)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening...")
	if err := http.ListenAndServe(":8080", &updater{client}); err != nil {
		log.Println(err)
	}
}

