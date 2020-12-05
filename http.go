package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func httpHandler(u Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		hostname := q.Get("hostname")
		newIP := strings.TrimSpace(q.Get("myip"))

		log.Println("Update DNS", hostname, "to", newIP)

		if err := u.UpdateDNS(newIP, hostname); err != nil {
			log.Println(err)

			var hpErr *HostnamePatternError
			if errors.As(err, &hpErr) {
				http.Error(w, "notfqdn", 400)
				return
			}

			var dnsErr *DNSProviderError
			if errors.As(err, &dnsErr) {
				http.Error(w, "dnserr", 500)
				return
			}
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte("good"))
	}
}
