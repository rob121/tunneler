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

Changes to `~/.tunnels` are picked up automatically. Changes to `~/.ssh/config` apply the next time a tunnel is started.

## Run

macOS opens bare executables from Finder in **Terminal**. Use the `.app` bundle instead:

```bash
make run
```

This builds `Tunneler.app` (menu bar only, no Terminal, no Dock icon) and launches it.

To install for everyday use (Finder, Spotlight, Login Items):

```bash
make install
```

Then launch **Tunneler** from `/Applications`.

For development from a terminal:

```bash
go run ./cmd/tunneler
```

When launched as an app, logs are written to `~/Library/Logs/Tunneler.log`.

## Menu

- Click a tunnel to start or stop it
- **Manage** — open `~/.tunnels` in your default editor
- **Start at Login** / **Quit** — provided by the menu bar framework

## Build

```bash
make          # creates Tunneler.app in this directory
make install  # copies Tunneler.app to /Applications
```

Do **not** double-click a bare binary built with `go build -o tunneler` — macOS will open Terminal. Always use `Tunneler.app`.
