package command

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hugolgst/rich-go/client"
)

type Serve struct {
	AppID   int    `required:"" env:"DISCORD_APP_ID" help:"The Discord Application/Client ID you'd like to use"`
	Bind    string `default:":1992" help:"Address and port to bind to"`
	Verbose bool   `default:"false" help:"Log each message handled"`
}

func (c *Serve) Run() error {
	if err := client.Login(strconv.Itoa(c.AppID)); err != nil {
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
				log.Printf("error handling message: %s", err)
				continue
			}
		}
	}()

	return <-done
}

type AugmentedActivity struct {
	*client.Activity
	Since string
}

func (c *Serve) handle(msg []byte) error {
	if c.Verbose {
		log.Printf("%s", msg)
	}

	activity := &AugmentedActivity{}

	if err := json.Unmarshal(msg, activity); err != nil {
		return err
	}

	switch since := activity.Since; since {
	case "":
	case "never":
	case "now":
		t := time.Now().Add(-1 * time.Second)
		activity.Timestamps = &client.Timestamps{Start: &t}
	default:
		secs, err := strconv.ParseInt(since, 10, 64)
		if err != nil {
			log.Printf("since: unparsable as int64; defaulting to never")
			activity.Timestamps = nil
			break
		}

		t := time.Unix(secs, 0)
		activity.Timestamps = &client.Timestamps{Start: &t}
	}

	return client.SetActivity(*activity.Activity)
}
