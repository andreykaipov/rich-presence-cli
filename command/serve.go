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
	cache   map[string]*client.Activity
}

func (c *Serve) BeforeApply() error {
	c.cache = map[string]*client.Activity{}
	return nil
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
	Since      string
	CacheKey   string
	CacheWrite string
}

func (c *Serve) handle(msg []byte) error {
	if c.Verbose {
		log.Printf("%s", msg)
	}

	activity := AugmentedActivity{}
	if err := json.Unmarshal(msg, &activity); err != nil {
		return err
	}

	c.resolveSince(&activity)
	c.resolveCache(&activity)

	return client.SetActivity(*activity.Activity)
}

// Since improves the UX for specifying activity timestamps
//
func (c *Serve) resolveSince(activity *AugmentedActivity) {
	var t time.Time

	switch since := activity.Since; since {
	case "":
	case "never":
	case "cached":
		// handled by cache resolution
	case "now":
		t = time.Now().Add(-1 * time.Second)
	default:
		secs, err := strconv.ParseInt(since, 10, 64)
		if err != nil {
			log.Printf("since: unparsable as int64; defaulting to never")
			break
		}

		t = time.Unix(secs, 0)
	}

	if t.IsZero() {
		activity.Timestamps = nil
		return
	}

	activity.Timestamps = &client.Timestamps{Start: &t}
}

// update given Activity with any cached values based on the cache key
//
func (c *Serve) resolveCache(activity *AugmentedActivity) {
	cached, present := c.cache[activity.CacheKey]
	if !present {
		goto write
	}

	if activity.Details == "cached" {
		activity.Details = cached.Details
	}

	if activity.State == "cached" {
		activity.State = cached.State
	}

	if activity.LargeImage == "cached" {
		activity.LargeImage = cached.LargeImage
	}

	if activity.LargeText == "cached" {
		activity.LargeText = cached.LargeText
	}

	if activity.SmallImage == "cached" {
		activity.SmallImage = cached.SmallImage
	}

	if activity.SmallText == "cached" {
		activity.SmallText = cached.SmallText
	}

	if activity.Since == "cached" {
		activity.Timestamps = cached.Timestamps
	}

write:
	if activity.CacheKey == "" {
		return
	}

	if activity.CacheWrite == "no" {
		return
	}

	if activity.CacheWrite == "if_not_present" {
		if !present {
			c.cache[activity.CacheKey] = activity.Activity
		}
		return
	}

	if activity.CacheWrite == "always" {
		c.cache[activity.CacheKey] = activity.Activity
	}
}
