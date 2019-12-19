package main

import (
	"github.com/kubatek94/dyndns-cloudflare-adapter/cf"
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

type updater struct {
	*cf.Client
}

func (u *updater) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var records []cf.DNSRecord

	q := r.URL.Query()
	hostname := q.Get("hostname")
	newIP := strings.TrimSpace(q.Get("myip"))
	log.Println("Update DNS ", hostname, " to ", newIP)

	if hostname != "" {
		pattern, err := regexp.Compile(hostname)
		if err != nil {
			log.Println(err)
			http.Error(w, "notfqdn", 400)
			return
		}

		records, err = u.FindDNSRecords(pattern)
	} else {
		records, err = u.FindDNSRecords(nil)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "nohost", 412)
		return
	}

	for _, record := range records {
		if record.IP != newIP {
			if err = u.UpdateDNSRecord(record, newIP); err != nil {
				log.Println(err)
				http.Error(w, "dnserr", 500)
				return
			}
		}
	}

	_, _ = w.Write([]byte("good"))
	w.WriteHeader(200)
}
