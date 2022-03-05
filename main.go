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

func main() {
	bind := flag.String("bind", ":1992", "Address and port to bind to")
	flag.Parse()

	appid := os.Getenv("DISCORD_APP_ID")
	if appid == "" {
		log.Fatal("DISCORD_APP_ID is empty")
	}

	if err := client.Login(appid); err != nil {
		log.Fatal(err)
	}

	log.Fatal(serve(*bind))
}

func serve(address string) error {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer pc.Close()

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
	activity := &client.Activity{}
	if err := json.Unmarshal(msg, activity); err != nil {
		return err
	}

	return client.SetActivity(*activity)
}
