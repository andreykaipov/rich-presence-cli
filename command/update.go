package command

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/alecthomas/kong"
	discord "github.com/hugolgst/rich-go/client"
)

type Update struct {
	Server     string            `default:":1992" help:"Address and port of the server"`
	Details    string            `help:"First line of your presence"`
	State      string            `help:"Second line of your presence"`
	LargeImage string            `help:"ID of the large asset for the activity"`
	LargeText  string            `help:"Text displayed when hovering over the large image"`
	SmallImage string            `help:"ID of the small asset for the activity"`
	SmallText  string            `help:"Text displayed when hovering over the small image"`
	Buttons    map[string]string `help:"Any buttons you might want, e.g. label=url"`
	Since      string            `default:"never" placeholder:"now|never|cached|<seconds-since-epoch>" help:"Time since the activity began; defaults to 'never'"`
	CacheKey   string            `help:"Key to read a previous activity update from the server cache, accessed by setting an arg to 'cached'"`
	CacheWrite string            `default:"no" placeholder:"no|if_not_present|always" help:"Write the current activity update to the cache with the provided --cache-key according to the provided write strategy; defaults to 'no'"`
	Dry        DryFlag           `help:"Dry run (prints the JSON payload to stdout)"`
}

type DryFlag string

func (f DryFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (f DryFlag) IsBool() bool                         { return true }
func (f DryFlag) AfterApply(app *kong.Kong, payload string) error {
	fmt.Println(payload)
	app.Exit(0)
	return nil
}

func (c *Update) AfterApply(ctx *kong.Context) error {
	buttons := []*discord.Button{}
	for label, url := range c.Buttons {
		buttons = append(buttons, &discord.Button{Label: label, Url: url})
	}

	activity := &AugmentedActivity{
		Activity: &discord.Activity{
			Details:    c.Details,
			State:      c.State,
			LargeImage: c.LargeImage,
			LargeText:  c.LargeText,
			SmallImage: c.SmallImage,
			SmallText:  c.SmallText,
			Buttons:    buttons,
		},
		Since:      c.Since,
		CacheKey:   c.CacheKey,
		CacheWrite: c.CacheWrite,
	}

	bytes, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	ctx.Bind(string(bytes))

	return nil
}

func (c *Update) Run(payload string) error {
	conn, err := net.Dial("udp", c.Server)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = fmt.Fprint(conn, payload); err != nil {
		return err
	}

	return nil
}
