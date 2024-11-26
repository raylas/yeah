# yeah

An API and CLI to return the assigned vendor for given MAC addresses.

## Sources

Queries are ran against the IEEE OUI vendor lists:

- <https://standards-oui.ieee.org/oui/oui.csv>
- <https://standards-oui.ieee.org/cid/cid.csv>
- <https://standards-oui.ieee.org/iab/iab.csv>
- <http://standards-oui.ieee.org/oui28/mam.csv>
- <https://standards-oui.ieee.org/oui36/oui36.csv>

## Deployment

```bash
just deploy
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
