# Tunneler

A lightweight macOS menu bar application for managing SSH SOCKS tunnels.

## Requirements

- macOS
- Go 1.24+
- SSH access configured in `~/.ssh/config`

## Philosophy

- Uses the system `/usr/bin/ssh`
- Uses your existing `~/.ssh/config`, keys, agent, and ProxyJump settings
- Stores only a simple list of tunnels in `~/.tunnels`

## Configuration

Create `~/.tunnels`:

```yaml
tunnels:
  - name: Example SOCKS
    ssh_host: example-host
    dynamic: 1080
```

Each tunnel references an SSH config host alias (`ssh_host`). Tunneler does not store hostnames, users, or keys.

Example `~/.ssh/config`:

```
Host example-host
    HostName example.com
    User alice
    IdentityFile ~/.ssh/id_ed25519
```

Toggling the tunnel in the menu bar runs:

```bash
ssh -N -D 1080 example-host
```

## Run

```bash
go run ./cmd/tunneler
```

The app appears as an icon in the menu bar. Use **Manage** in the menu to open `~/.tunnels` in your default editor.

## Build

```bash
go build -o tunneler ./cmd/tunneler
```
