# yeah

An API and CLI to return the assigned vendor for given MAC addresses.

## Sources

Queries are ran against the IEEE OUI vendor lists:

- <https://standards-oui.ieee.org/oui/oui.csv>
- <https://standards-oui.ieee.org/cid/cid.csv>
- <https://standards-oui.ieee.org/iab/iab.csv>
- <http://standards-oui.ieee.org/oui28/mam.csv>
- <https://standards-oui.ieee.org/oui36/oui36.csv>

## Usage

### Data

```bash
go run cmd/data/main.go
```

### CLI

Install:

- [Binary releases](https://github.com/raylas/yeah/releases)

```bash
Usage: yeah [-w] [-o OUTPUT] [-l] [-b BIND] [-v LOGLEVEL] [MACS [MACS ...]]

Positional arguments:
  MACS

Options:
  -w                     include additional fields
  -o OUTPUT              output format (table,json,html) [default: table]
  -l                     run server
  -b BIND                server bind address [default: :8080]
  -v LOGLEVEL            log level (info,debug) [default: info]
  --help, -h             display this help and exit
```

## Deployment

```bash
flyctl deploy --remote-only
```

### Initial

Create Fly app:

```bash
fly apps create --name yeah
```

Secrets:

```bash
fly secrets set OTEL_EXPORTER_OTLP_HEADERS=x-honeycomb-team=<api_key>
```

Create a token and place as a GitHub repository secret named `FLY_API_TOKEN`:

```bash
fly tokens create deploy --app yeah --expiry 999999h
```

### DNS

Once the application is deployed, grab its IPv6 and IPv4 addresses and create a
pair of A/AAAA records for your domain.

Then run:

```bash
fly certs create <domain_name>
```

If behind Cloudflare's proxy, it'll prompt you to add a `CNAME` record, too.
