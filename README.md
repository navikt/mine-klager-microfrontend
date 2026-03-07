# Mine klager - Min side microfrontend

A microfrontend for [Min side](https://www.nav.no/minside) that links users to [Mine klager](https://mine-klager.nav.no).

A build-time TypeScript script renders the [LinkCard](https://aksel.nav.no/komponenter/core/linkcard) component to a static HTML template, and a Go code generator reads it along with `@navikt/ds-css` files, tree-shakes design tokens, discovers needed component CSS, and outputs a single CSS file. A minimal Go server embeds the template and CSS, assembles per-language HTML fragments at startup, and serves them.

Each fragment is wrapped in a [Declarative Shadow DOM](https://caniuse.com/declarative-shadow-dom) (`<template shadowrootmode="open">`) for style isolation.

> **Since February 2024, this feature works across the latest devices and browser versions.**

- MDN: [Declarative Shadow DOM](https://developer.mozilla.org/en-US/docs/Web/API/Web_components/Using_shadow_DOM#declaratively_with_html)
- MDN: [shadowRootMode](https://developer.mozilla.org/en-US/docs/Web/API/HTMLTemplateElement/shadowRootMode)
- MDN: [Shadow DOM](https://web.dev/articles/declarative-shadow-dom)
- Can I Use: [Declarative Shadow DOM](https://caniuse.com/declarative-shadow-dom)
- web.dev: [Declarative Shadow DOM](https://web.dev/articles/declarative-shadow-dom)

Uses [LinkCard](https://aksel.nav.no/komponenter/core/linkcard) from `@navikt/ds-react` with tree-shaken design tokens from `@navikt/ds-css`.

## How it works

1. **Build time:** `bun run template` renders the `LinkCard` component from `@navikt/ds-react` with `renderToStaticMarkup` and writes `templates/template.html`. Then `go run ./generate/css` reads that template, extracts CSS classes directly from it, discovers needed component CSS from `@navikt/ds-css`, tree-shakes tokens, and writes `templates/style.css`.
2. **Startup:** The Go server (`main.go`) embeds `templates/template.html` and `templates/style.css` at compile time (`//go:embed`), then assembles per-language HTML fragments by replacing `{{TITLE}}`, `{{DESCRIPTION}}`, and `{{URL}}` placeholders.
3. Each HTML fragment uses Declarative Shadow DOM, so styles are fully encapsulated and won't conflict with the host page.

### CSS pipeline (build time)

1. **Component CSS discovery** — Extracts CSS class names from the template HTML and matches them against `@navikt/ds-css/dist/component/*.min.css` files. Only CSS files with overlapping selectors are included.
2. **Token tree-shaking** — Parses `tokens.css` to build a dependency graph of custom properties, walks it from the properties referenced in the consumer CSS, and strips everything else.
3. **Shadow DOM scoping** — Rewrites `:root` to `:host` so custom properties are scoped to the shadow boundary.

## Routes

- `/nb` — Norwegian Bokmål
- `/nn` — Norwegian Nynorsk
- `/en` — English
- `/isAlive` — Liveness probe
- `/isReady` — Readiness probe

## Environment Variables

- `NAIS_CLUSTER_NAME` — Determines which domain to use for links:
  - `prod-gcp` → `https://mine-klager.nav.no`
  - Any other value → `https://mine-klager.ansatt.dev.nav.no`
- `PORT` — Server port (default: `8080`)

## Local Development

```sh
bun i
bun run template
go run ./generate/css
go run .
```

**Handy one-liner for local development**

```sh
bun run template && go run ./generate/css && go run .
```

### Docker Compose

There is also a `docker-compose.yml` for local development that runs the Go server in a container.

```sh
bun i
bun run template
go run ./generate/css
CGO_ENABLED=0 go build -o server .
docker compose up --build
```

Then open:
- [localhost:8080/nb](http://localhost:8080/nb)
- [localhost:8080/nn](http://localhost:8080/nn)
- [localhost:8080/en](http://localhost:8080/en)

## Testing

```sh
go test ./...
```

## Why Bun?

Bun is used to install npm packages from GitHub's npm registry and to run the TypeScript template generator (`generate/template/index.tsx`), which renders the `LinkCard` React component to static HTML. All CSS processing and serving is done in Go.

## Project Structure

```
.
├── generate/
│   ├── css/
│   │   ├── main.go             # Build-time CSS generator
│   │   ├── build.go            # Assembles final CSS (global + component + tokens)
│   │   ├── components.go       # Discovers component CSS files from HTML class usage
│   │   ├── components_test.go  # Tests for component CSS discovery
│   │   ├── treeshake.go        # Removes unused design tokens from CSS
│   │   ├── treeshake_test.go   # Tests for token tree-shaking
│   │   └── format-bytes.go     # Byte formatting utility
│   └── template/
│       ├── index.tsx            # Renders LinkCard to static HTML template
│       └── icon.tsx             # Custom illustration icon component
├── templates/
│   ├── template.html           # Generated HTML template (build output, gitignored)
│   └── style.css               # Generated CSS (build output, gitignored)
├── main.go                     # Go HTTP server (embeds templates/, assembles HTML at startup)
├── go.mod                      # Go module
├── package.json                # npm dependencies and build script
├── bun.lock                    # Bun lockfile
├── bunfig.toml                 # Bun configuration
├── tsconfig.json               # TypeScript configuration
├── biome.json                  # Biome linter and formatter configuration
├── .tool-versions              # Tool version manager (Bun version)
├── Dockerfile                  # Minimal scratch-based Docker image
├── docker-compose.yml          # Local development with Docker
├── LICENSE                     # MIT License
├── CODEOWNERS                  # GitHub code owners
├── nais/
│   ├── dev.yaml                # NAIS config for dev-gcp
│   └── prod.yaml               # NAIS config for prod-gcp
└── .github/
    └── workflows/
        ├── build.yaml          # Reusable build workflow
        ├── build-and-deploy.yaml
        └── deploy-dev.yaml
```
