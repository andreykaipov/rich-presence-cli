## rich-presence-udp

It's a UDP proxy server that allows you to update your Discord Rich Presence.

## usage

Run it:

```console
❯ make serve
go build -o rich-presence.exe main.go
./rich-presence.exe serve --verbose
2022/03/26 17:37:09 Listening on :1992
```

In another shell, send an update:

```console
❯ ./rich-presence.exe update --details "$(basename "$PWD")" --state browsing...
```

In the above snippets, I'm running on Windows via WSL (hence the `.exe` suffix),
but it should work just fine on desktop Linux too.

We can also just form the JSON payload ourselves:

```console
❯ echo '{"details":"wassup"}' > /dev/udp/$WINHOST/1992
```

### additional configuration

Check out the [`rich-presence.yml`](./rich-presence.yml) configuration file at
the root of this repo to see everything that's configurable. The search paths
for this file are as follows:

- `rich-presence.yml`
- `~/.rich-presence.yml`
- `~/.config/rich-presence/rich-presence.yml`

## example

Go crazy:

![todo gif](./example.gif)
