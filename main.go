package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubatek94/dyndns-cloudflare-adapter/cf"
)

type modeHandler func(context.Context, Updater) error

func main() {
	token := os.Getenv("CF_API_TOKEN")

	if token == "" {
		log.Fatalln("CF_API_TOKEN missing from environment")
	}

	client, err := cf.NewClient(token)
	if err != nil {
		log.Fatalln(err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	stop, errs := execHandler(Updater{client}, pickModeHandler())

	select {
	case <-sigint:
		stop()
	case err := <-errs:
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Printf("usage: %s <http|stun> [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func pickModeHandler() modeHandler {
	httpFlagSet := flag.NewFlagSet("http", flag.ExitOnError)
	port := httpFlagSet.String("port", "8080", "port of the http server")

	stunFlagSet := flag.NewFlagSet("stun", flag.ExitOnError)
	hp := stunFlagSet.String("hp", ".+", "hostname pattern to match DNS records for update")

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "http":
		httpFlagSet.Parse(os.Args[2:])
		return func(ctx context.Context, u Updater) error {
			return httpMode(ctx, u, *port)
		}
	case "stun":
		stunFlagSet.Parse(os.Args[2:])
		return func(ctx context.Context, u Updater) error {
			return stunMode(ctx, u, *hp)
		}
	default:
		usage()
	}

	return nil
}

func execHandler(u Updater, handler modeHandler) (func(), <-chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error)

	go func() {
		errs <- handler(ctx, u)
	}()

	return func() {
		cancel()
		<-errs
	}, errs
}

func httpMode(ctx context.Context, u Updater, port string) error {
	srv := http.Server{
		Addr:    ":" + port,
		Handler: httpHandler(u),
	}

	errs := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-errs:
		return err
	}

	return nil
}

func stunMode(ctx context.Context, u Updater, hp string) error {
	stunUpdate(u, hp)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			stunUpdate(u, hp)
		}
	}
}
