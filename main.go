package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/andreykaipov/rich-presence-udp/command"
	"github.com/goccy/go-yaml"
)

var cli struct {
	Config kong.ConfigFlag `type:"path" help:"Path to a YAML file with defaults"`
	Serve  command.Serve   `cmd:"" help:"Start a UDP proxy server for Discord IPC"`
	Update command.Update  `cmd:"" help:"Send an update to the server"`
}

func main() {
	home := os.Getenv("HOME") // WSL workaround; invoke with WSLENV=HOME/w

	ctx := kong.Parse(
		&cli,
		kong.Name("rich-presence-udp"),
		kong.Description("Manage your Discord Rich Presence via UDP"),
		kong.Configuration(
			yamlEnvResolver,
			"rich-presence.yml",
			fmt.Sprintf("%s/.rich-presence.yml", home),
			fmt.Sprintf("%s/.config/rich-presence/rich-presence.yml", home),
		),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
	)

	ctx.FatalIfErrorf(ctx.Run())
}

// resolves configs to init our commands' flags
func yamlEnvResolver(r io.Reader) (kong.Resolver, error) {
	values := map[string]interface{}{}

	if err := yaml.NewDecoder(r).Decode(&values); err != nil {
		return nil, err
	}

	var f kong.ResolverFunc = func(context *kong.Context, _ *kong.Path, flag *kong.Flag) (interface{}, error) {
		name := strings.ReplaceAll(flag.Name, "-", "_")
		val, ok := values[name]
		if !ok {
			return nil, nil
		}

		switch v := val.(type) {
		case string:
			k := "envexpand:"
			if strings.HasPrefix(v, k) {
				val = os.ExpandEnv(v[len(k):])
			}
		}

		return val, nil
	}

	return f, nil
}
