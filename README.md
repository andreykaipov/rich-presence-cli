## rich-presence-udp

It's a UDP server that allows you to update your Discord Rich Presence.

## usage

Run it:

```console
❯ make run
go build -o rich-presence.exe main.go
./rich-presence.exe -verbose
2022/03/05 10:39:49 Listening on :1992
```

In another shell, send a UDP message to it:

```console
❯ jq -cM . example.json > /dev/udp/$WINHOST/1992
```

In the above, I'm running this on Windows via WSL, but it should work just fine
on desktop Linux too.

## example

Go crazy:

![gif of automated rich presence updates via a shell loop](./example.gif)
