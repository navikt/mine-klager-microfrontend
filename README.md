# Mine klager - Min side microfrontend

Go-based microfrontend with HTML templating and environment-aware domain selection.

The project uses CSS nesting, which [is widely supported](https://caniuse.com/css-nesting).

[Nav only requires support for the latest major version of Chrome, Firefox, Edge and Safari.](https://nav-it.slack.com/archives/C04V21LT27P/p1767786403388759?thread_ts=1766063231.776659&cid=C04V21LT27P)

All CSS classes are prefixed with `e49b256-` to avoid conflicts with other CSS on Min Side.

## Routes

- `/nb` - Norwegian Bokmål
- `/nn` - Norwegian Nynorsk
- `/en` - English
- `/fallback` - Norwegian Bokmål
- `/isAlive` - Liveness probe
- `/isReady` - Readiness probe

## Environment Variables

- `NAIS_CLUSTER_NAME` - Determines which domain to use in URLs:
  - `prod-gcp` → `https://mine-klager.nav.no`
  - `dev-gcp` (or any other value) → `https://mine-klager.ansatt.dev.nav.no`
- `PORT` - Server port (default: `8080`)

## Local Development

### With Docker Compose

```sh
CGO_ENABLED=0 go build -o server .
docker compose up --build
```

### With Docker

```sh
CGO_ENABLED=0 go build -o server .
docker build -t mine-klager-microfrontend .
docker run -p 8080:8080 mine-klager-microfrontend
```

### With Go

```sh
go run .
```

Then open:
- [localhost:8080/nb](http://localhost:8080/nb)
- [localhost:8080/nn](http://localhost:8080/nn)
- [localhost:8080/en](http://localhost:8080/en)
- [localhost:8080/fallback](http://localhost:8080/fallback)

## Project Structure

```
.
├── main.go                 # Go HTTP server
├── templates/
│   ├── style.css           # Embedded CSS styles
│   └── template.html       # Go HTML template
├── go.mod                  # Go module file
├── Dockerfile              # Docker image definition
├── nais/
│   ├── dev.yaml            # NAIS config for dev-gcp
│   └── prod.yaml           # NAIS config for prod-gcp
└── docker-compose.yml      # Local development
```
