package main

import (
	"github.com/ccding/go-stun/stun"
	"github.com/kubatek94/dyndns-cloudflare-adapter/cf"
	
	// "log"
	// "net/http"
	// "regexp"
	// "strings"
	// "errors"
	// "fmt"
)

type updater struct {
	cf *cf.Client
}

// var hostnameInvalidError = errors.New("Hostname pattern is not valid.")


func (u *updater) UpdateDNS(newIP string, hostnamePattern string) error {

	return nil

	// var err error
	// var records []cf.DNSRecord

	// if hostnamePattern != "" {
	// 	pattern, err := regexp.Compile(hostnamePattern)
	// 	if err != nil {
	// 		fmt.Errorf("%w: hostname pattern invalid", err)

	// 		log.Println(err)
	// 		http.Error(w, "notfqdn", 400)
	// 		return
	// 	}

	// 	records, err = u.cf.FindDNSRecords(pattern)
	// } else {
	// 	records, err = u.cf.FindDNSRecords(nil)
	// }

	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "nohost", 412)
	// 	return
	// }

	// for _, record := range records {
	// 	if record.IP != newIP {
	// 		if err = u.UpdateDNSRecord(record, newIP); err != nil {
	// 			log.Println(err)
	// 			http.Error(w, "dnserr", 500)
	// 			return
	// 		}
	// 	}
	// }
}

// func (u *updater) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	var records []cf.DNSRecord

// 	q := r.URL.Query()
// 	hostname := q.Get("hostname")
// 	newIP := strings.TrimSpace(q.Get("myip"))
// 	log.Println("Update DNS ", hostname, " to ", newIP)

// 	if hostname != "" {
// 		pattern, err := regexp.Compile(hostname)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "notfqdn", 400)
// 			return
// 		}

// 		records, err = u.FindDNSRecords(pattern)
// 	} else {
// 		records, err = u.FindDNSRecords(nil)
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "nohost", 412)
// 		return
// 	}

// 	for _, record := range records {
// 		if record.IP != newIP {
// 			if err = u.UpdateDNSRecord(record, newIP); err != nil {
// 				log.Println(err)
// 				http.Error(w, "dnserr", 500)
// 				return
// 			}
// 		}
// 	}

// 	w.WriteHeader(200)
// 	_, _ = w.Write([]byte("good"))
// }
