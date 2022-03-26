package main

import (
	"github.com/alecthomas/kong"
	"github.com/andreykaipov/rich-presence-udp/command"
)

var cli struct {
	Serve  command.Serve  `cmd:"" help:"Start the server"`
	Update command.Update `cmd:"" help:"Send an update to the server"`
}

func main() {
	ctx := kong.Parse(
		&cli,
		kong.Name("rich-presence-udp"),
		kong.Description("Manage your Discord Rich Presence via UDP"),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
	)

	ctx.FatalIfErrorf(ctx.Run())
}
