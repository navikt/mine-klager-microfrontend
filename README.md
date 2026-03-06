# Mine klager - Min side microfrontend

A microfrontend for [Min side](https://www.nav.no/minside) that links users to [Mine klager](https://mine-klager.nav.no).

Built with React SSR (via Bun) at build time, served as static HTML fragments by a Go server.

Each fragment is wrapped in a [Declarative Shadow DOM](https://caniuse.com/declarative-shadow-dom) (`<template shadowrootmode="open">`) for style isolation.

> **Since February 2024, this feature works across the latest devices and browser versions.**

- MDN: [Declarative Shadow DOM](https://developer.mozilla.org/en-US/docs/Web/API/Web_components/Using_shadow_DOM#declaratively_with_html)
- MDN: [shadowRootMode](https://developer.mozilla.org/en-US/docs/Web/API/HTMLTemplateElement/shadowRootMode)
- MDN: [Shadow DOM](https://web.dev/articles/declarative-shadow-dom)
- Can I Use: [Declarative Shadow DOM](https://caniuse.com/declarative-shadow-dom)
- web.dev: [Declarative Shadow DOM](https://web.dev/articles/declarative-shadow-dom)

Uses [LinkCard](https://aksel.nav.no/komponenter/core/linkcard) from `@navikt/ds-react` with tree-shaken design tokens from `@navikt/ds-css`.

## How it works

1. `bun run build` renders React components to static HTML with `react-dom/server`, tree-shakes `@navikt/ds-css` tokens, and outputs self-contained `.html` fragments into `dist/`.
2. The Go server embeds the `dist/` directory at compile time (`//go:embed`), replaces `{{BASE_URL}}` with the environment-appropriate domain, and serves the fragments.
3. Each HTML fragment uses Declarative Shadow DOM, so styles are fully encapsulated and won't conflict with the host page.

## Routes

- `/nb` - Norwegian Bokmål
- `/nn` - Norwegian Nynorsk
- `/en` - English
- `/isAlive` - Liveness probe
- `/isReady` - Readiness probe

## Environment Variables

- `NAIS_CLUSTER_NAME` - Determines which domain to use for links:
  - `prod-gcp` → `https://mine-klager.nav.no`
  - Any other value → `https://mine-klager.ansatt.dev.nav.no`
- `PORT` - Server port (default: `8080`)

## Local Development

```sh
bun i
bun run build
go run .
```

### Docker Compose

There is also a `docker-compose.yml` for local development that runs the Go server in a container.
This allows you to test the production-like environment locally, including the correct handling of environment variables.

```sh
bun i
bun run build
CGO_ENABLED=0 go build -o server .
docker compose up --build
```

Then open:
- [localhost:8080/nb](http://localhost:8080/nb)
- [localhost:8080/nn](http://localhost:8080/nn)
- [localhost:8080/en](http://localhost:8080/en)

## Scripts

| Script          | Description                                 |
| --------------- | ------------------------------------------- |
| `bun run build` | Generate static HTML fragments into `dist/` |
| `bun test`      | Run tests                                   |
| `bun typecheck` | Run TypeScript type checking                |
| `bun lint`      | Run Biome linter                            |

## Project Structure

```
.
├── src/
│   ├── generate.tsx            # Build script — renders React to static HTML
│   ├── microfrontend.tsx       # Main React component (LinkCard)
│   ├── icon.tsx                # SVG icon component
│   ├── css.ts                  # Builds CSS for Shadow DOM (global + auto-discovered components)
│   ├── find-component-css.ts   # Discovers component CSS files from HTML class usage
│   ├── find-component-css.test.ts
│   ├── tree-shake-tokens.ts    # Removes unused design tokens from CSS
│   ├── tree-shake-tokens.test.ts
│   └── format.ts               # Byte formatting utility
├── dist/                       # Generated HTML fragments (build output)
│   ├── nb.html
│   ├── nn.html
│   └── en.html
├── main.go                     # Go HTTP server (embeds dist/)
├── go.mod                      # Go module
├── package.json                # Bun/Node dependencies and scripts
├── bun.lock                    # Bun lockfile
├── tsconfig.json               # TypeScript config
├── biome.json                  # Biome linter/formatter config
├── Dockerfile                  # Minimal scratch-based Docker image
├── docker-compose.yml          # Local development with Docker
├── nais/
│   ├── dev.yaml                # NAIS config for dev-gcp
│   └── prod.yaml               # NAIS config for prod-gcp
└── .github/
    └── workflows/
        └── build-and-deploy.yaml  # CI/CD pipeline
```
