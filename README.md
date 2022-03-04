## rich-presence-udp

It's a UDP server that allows you to update your Discord Rich Presence.

## usage

### Windows via WSL

After invoking the executable for the first time, you'll have to accept the
Windows Firewall security alert to allow inbound UDP traffic to your Windows
host.

```console
❯ GOOS=windows go build -o rich-presence main.go
❯ DISCORD_APP_ID=942604338927374438 WSLENV=DISCORD_APP_ID/w ./rich-presence.exe
❯ echo '{"state":"hello"}' | nc -uw0 "$WSL_HOST" 1992
```

### Linux

It hasn't really been tested on Linux desktop, but I imagine it'd work just fine
since the socket [rich-go](https://github.com/hugolgst/rich-go) expects (e.g.
`$XDG_RUNTIME_DIR/discord-ipc-0`) would actually exist, if you're running
Discord for Linux.

```console
❯ go build -o rich-presence main.go
❯ DISCORD_APP_ID=942604338927374438 ./rich-presence
```

## misc

Some notes on the Windows via WSL invocation:

- Assuming Discord is running on the Windows host, and not within WSL, we must
  build the binary with `GOOS=windows` so the underlying
  [rich-go](https://github.com/hugolgst/rich-go) library can read from the
  named pipe Discord uses on Windows (i.e. `\\.\\pipe\\discord-ipc-0`).

  We can check for the existence of this named pipe as follows:

  ```console
  ❯ powershell.exe '[System.IO.Directory]::GetFiles("\\.\\pipe\\")' | grep discord
  ```

- Since we're invoking Windows executables from within WSL, we must also tell
  WSL to share with Windows any env vars we've set using
  `WSLENV=DISCORD_APP_ID/w` (see [this blog
  post](https://devblogs.microsoft.com/commandline/share-environment-vars-between-wsl-and-windows/)
  for more details).

- After the server is running, we need a way to send messages to it. However,
  since we invoked a Windows executable, the server is bound to the host of the
  WSL instance. We can find that as follows:

  ```console
  ❯ powershell.exe "
          Get-NetAdapter 'vEthernet (WSL)' |
          Get-NetIPAddress -AddressFamily IPv4 |
          Select -Expand IPAddress
  "
  ```

  Note this is rather slow. I'd recommend adding this to the environment vars
  your login shell sets. I have mine set under `WSL_HOST`, like in the example
  above.
