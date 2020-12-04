package main

import "github.com/ccding/go-stun/stun"

func run() {
	nat, host, err := stun.NewClient().Discover()
}
