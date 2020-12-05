package main

import (
	"github.com/ccding/go-stun/stun"
	"log"
)

const stunServer = "stun.l.google.com:19302"

func stunUpdate(u Updater, hostnamePattern string) {
	client := stun.NewClient()
	client.SetServerAddr(stunServer)

	_, host, err := client.Discover()
	if host == nil {
		log.Println("failed detecting our IP address:", err)
		return
	}

	log.Println("Update DNS", hostnamePattern, "to", host.IP())

	if err := u.UpdateDNS(host.IP(), hostnamePattern); err != nil {
		log.Println(err)
	}
}
