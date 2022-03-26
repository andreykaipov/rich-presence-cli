These are some miscellaneous notes for myself (or whoever else reads this) about
running this on Windows within WSL.

- Assuming Discord is running on the Windows host, and not within WSL, we
  also need to run our server on the Windows host so the underying
  [rich-go](https://github.com/hugolgst/rich-go) library can interface with the
  named pipe Discord uses on Windows (i.e. `\\.\\pipe\\discord-ipc-0`).

  We can check for the existence of this named pipe as follows:

  ```console
  ❯ powershell.exe '[System.IO.Directory]::GetFiles("\\.\\pipe\\")' | grep discord
  ```

- Since we're invoking a Windows executable from within WSL, we must also tell
  WSL to share with Windows any env vars we've set via
  `WSLENV=DISCORD_APP_ID/w:HOME/w`. See [this blog
  post](https://devblogs.microsoft.com/commandline/share-environment-vars-between-wsl-and-windows/)
  for more details.

- After invoking the executable for the first time, we'll have to accept the
  Windows Firewall security alert to allow inbound UDP traffic to our Windows
  host.

- Sending messages to our server won't work via `localhost` because the server
  is bound on the Windows host. We can talk to the Windows host using the
  default gateway within our WSL instance:

  ```console
  ❯ ip r show default | awk '{print $3}'
  172.22.128.1
  ```

  We could also ask Powershell:

  ```console
  ❯ powershell.exe '
          Get-NetIPAddress -AddressFamily IPv4 -InterfaceAlias "vEthernet (WSL)" |
          Select -ExpandProperty IPAddress
  '
  ```

- Once everything is up, we can confirm our binary is actually running
  within Windows:

  ```console
  ❯ powershell.exe 'Get-Process -Name rich-presence | Select -ExpandProperty Id'
  20464
  ```

  ```console
  ❯ powershell.exe 'Get-NetUDPEndpoint -OwningProcess 20464'

  LocalAddress                             LocalPort
  ------------                             ---------
  ::                                       1992
  ```
