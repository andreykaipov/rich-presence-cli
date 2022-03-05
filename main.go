package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"

	"github.com/hugolgst/rich-go/client"
)

var (
	bind    string
	verbose bool
)

func init() {
	flag.StringVar(&bind, "bind", ":1992", "Address and port to bind to")
	flag.BoolVar(&verbose, "verbose", false, "Log each message received")
	flag.Parse()
}

func main() {
	appid := os.Getenv("DISCORD_APP_ID")
	if appid == "" {
		log.Fatal("DISCORD_APP_ID is empty")
	}

	if err := client.Login(appid); err != nil {
		log.Fatal(err)
	}

	log.Fatal(serve(bind))
}

func serve(address string) error {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer pc.Close()

	log.Printf("Listening on %s", address)

	done := make(chan error, 1)
	buffer := make([]byte, 1024)

	go func() {
		for {
			n, _, err := pc.ReadFrom(buffer)
			if err != nil {
				done <- err
				return
			}

			msg := bytes.TrimRight(buffer[:n], "\n")

			if err := handle(msg); err != nil {
				log.Printf("Error handling message %q. Ignoring...", msg)
				continue
			}
		}
	}()

	return <-done
}

func handle(msg []byte) error {
	if verbose {
		log.Printf("%s", msg)
	}

	activity := &client.Activity{}

	if err := json.Unmarshal(msg, activity); err != nil {
		return err
	}

	return client.SetActivity(*activity)
}
