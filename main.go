package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/hugolgst/rich-go/client"
)

func main() {
	if err := client.Login(os.Getenv("DISCORD_APP_ID")); err != nil {
		log.Fatal(err)
	}

	log.Fatal(serve("0.0.0.0:1992"))
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
