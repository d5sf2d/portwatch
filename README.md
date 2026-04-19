# portwatch

> CLI tool to monitor and alert on open ports and service changes on a host

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start monitoring all open ports on the local host:

```bash
portwatch watch
```

Scan and display currently open ports:

```bash
portwatch scan
```

Monitor specific ports and alert on changes:

```bash
portwatch watch --ports 22,80,443 --interval 30s --alert email
```

Example output:

```
[2024-01-15 10:32:01] PORT OPENED: 8080/tcp (http-alt)
[2024-01-15 10:45:17] PORT CLOSED: 3000/tcp (ppp)
[2024-01-15 11:02:44] SERVICE CHANGED: 443/tcp nginx/1.18 -> nginx/1.24
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `60s` | Polling interval |
| `--ports` | all | Comma-separated list of ports to watch |
| `--alert` | `stdout` | Alert method: `stdout`, `email`, `webhook` |
| `--host` | `localhost` | Target host to monitor |

## License

MIT © [yourusername](https://github.com/yourusername)