The SSH server written in Go to log login attempts.
Password authentications always fail so no terminal access is given to the attacker.

For better testing experience add this to your `~/.ssh/config file`

```
Host localhost
    UserKnownHostsFile /dev/null
```
