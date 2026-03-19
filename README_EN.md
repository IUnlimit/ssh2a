# SSH2A

English | [简体中文](./README.md)

SSH2A is an SSH access authorization gateway that provides web-based authentication and honeypot protection for SSH ports.

## How It Works

```
Client --SSH--> [9022] SSH2A --Check IP state-->
  ├─ Verified   → Forward to local port 22 (normal SSH login)
  ├─ First deny → Record IP, close connection
  └─ Retry      → Enter honeypot (capture credentials, then disconnect)

Client --HTTP--> [9080] SSH2A Web UI
  └─ Enter password / 2FA → Verified → SSH access granted
```

## Features

- **Web Authentication** — Verify via browser with password or 2FA code, then SSH port is opened for that IP
- **API Authentication** — Supports `Authorization` header for programmatic access
- **SSH Honeypot** — Unverified IPs that retry SSH within the timeout window enter a honeypot that captures login credentials
- **Admin Panel** — Web dashboard for viewing captured credentials, rejected IPs, verified IPs and statistics
- **IP Whitelist** — Admin API endpoints can be restricted to specific IPs
- **Light/Dark Mode** — Frontend theme toggle with localStorage persistence
- **PostgreSQL Storage** — All records persisted to database
- **Single Binary** — Frontend embedded into Go binary via `embed`

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+ & pnpm
- PostgreSQL

### Build from Source

```bash
git clone https://github.com/IUnlimit/ssh2a.git
cd ssh2a

# Build frontend + compile binary
make all
```

Output binaries are in the `output/` directory.

### Docker Compose

```bash
# Build binary (Linux amd64)
make linux

# Start services
docker compose up -d
```

A default `config.yml` is generated on first run. Edit as needed and restart.

### Run Directly

```bash
./output/ssh2a_linux
```

A `config.yml` will be generated in the current directory on first run.

## Configuration

```yaml
bind:
  host: 0.0.0.0
  http-port: 9080        # Web UI port
  ssh-port: 9022         # SSH proxy port

authorization:
  type: 'basic'          # basic | authenticator
  basic:
    secret: '123456'
  authenticator:
    private-secret: ''   # 2FA private key

database:
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: '123456'
  dbname: ssh2a

honeypot:
  trigger-timeout: 3m    # Time before honeypot triggers after rejection

auth:
  valid-duration: 30m    # How long SSH access remains valid after verification

admin:
  allowed-hosts:         # Admin API IP whitelist, empty = no restriction
    - 127.0.0.1
    - ::1
```

## API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/auth` | Authenticate (JSON body or Authorization header) |
| GET | `/api/v1/status` | Query current IP status |
| GET | `/api/v1/admin/stats` | Overview statistics |
| GET | `/api/v1/admin/honeypot` | Honeypot captured credentials |
| GET | `/api/v1/admin/rejected` | Rejected IP list |
| GET | `/api/v1/admin/verified` | Verified IP list |

## License

[AGPL-3.0](./LICENSE)
