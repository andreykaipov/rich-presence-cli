package command

import (
	"bytes"
	"encoding/json"
	"log"
	"net"

	"github.com/alecthomas/kong"
	"github.com/hugolgst/rich-go/client"
)

type Serve struct {
	AppID   string `env:"DISCORD_APP_ID" help:"The Discord Application/Client ID you'd like to use"`
	Bind    string `default:":1992" help:"Address and port to bind to"`
	Verbose bool   `default:"false" help:"Log each message handled"`
}

func (c *Serve) Run(ctx *kong.Context) error {
	if err := client.Login(c.AppID); err != nil {
		return err
	}

	return c.start()
}

func (c *Serve) start() error {
	pc, err := net.ListenPacket("udp", c.Bind)
	if err != nil {
		return err
	}
	defer pc.Close()

	log.Printf("Listening on %s", c.Bind)

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

			if err := c.handle(msg); err != nil {
				log.Printf("Error handling message %s. Ignoring...", msg)
				continue
			}
		}
	}()

	return <-done
}

func (c *Serve) handle(msg []byte) error {
	if c.Verbose {
		log.Printf("%s", msg)
	}

	activity := &client.Activity{}

	if err := json.Unmarshal(msg, activity); err != nil {
		return err
	}

	return client.SetActivity(*activity)
}
