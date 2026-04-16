# portwatch

Lightweight CLI to monitor and alert on port changes across local and remote hosts.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

```bash
# Watch ports on localhost
portwatch --host localhost --ports 80,443,8080

# Monitor a remote host with alerts
portwatch --host 192.168.1.10 --ports 22,3306 --interval 30s --alert

# Watch multiple hosts from a config file
portwatch --config hosts.yaml
```

Example `hosts.yaml`:

```yaml
hosts:
  - address: localhost
    ports: [80, 443, 8080]
  - address: 192.168.1.10
    ports: [22, 3306, 5432]
interval: 30s
alert: true
```

portwatch will notify you whenever a port opens or closes, logging changes with timestamps to stdout or a specified log file.

```bash
[2024-01-15 10:32:01] OPEN   192.168.1.10:3306
[2024-01-15 10:35:44] CLOSED localhost:8080
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `localhost` | Target host to monitor |
| `--ports` | | Comma-separated list of ports |
| `--interval` | `60s` | Polling interval |
| `--config` | | Path to YAML config file |
| `--alert` | `false` | Enable desktop notifications |
| `--log` | | Path to log output file |

## License

MIT © 2024 [yourusername](https://github.com/yourusername)