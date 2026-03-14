# Contributing to godoceasy 📘

Thank you for considering contributing! Every contribution — bug fix, feature, doc improvement — makes Go documentation easier for beginners worldwide.

## Quick Start for Contributors

```bash
# 1. Fork and clone
git clone https://github.com/<your-username>/godoceasy.git
cd godoceasy

# 2. Build
go build -o godoceasy ./cmd/godoceasy

# 3. Test with any Go package
./godoceasy fmt
./godoceasy net/http

# 4. Make changes, rebuild, test
go build -o godoceasy ./cmd/godoceasy && go vet ./...
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     CLI (main.go)                       │
│  Parses args → Calls Fetcher → Parser → Generator →    │
│  Starts Web Server                                      │
└──────────┬──────────────────────────────────────────────┘
           │
           ▼
┌────────────────────┐   ┌────────────────────────┐
│  Fetcher           │   │  Parser                 │
│  (github.go)       │   │  (ast_parser.go)        │
│                    │   │                          │
│  • Stdlib → GOROOT │   │  • go/parser + go/ast   │
│  • Remote → git    │   │  • Extracts structs,    │
│    clone           │   │    interfaces, funcs,    │
│  • Vanity → known  │   │    methods, consts,      │
│    repos map       │   │    variables             │
│  • Version check   │   │  • 3-level recursive     │
└────────┬───────────┘   └────────┬─────────────────┘
         │                        │
         ▼                        ▼
┌─────────────────────────────────────────────────────────┐
│  Generator (docs/generator.go)                          │
│                                                         │
│  • Converts raw AST docs → beginner-friendly format     │
│  • Detects builder pattern (WithX / SetX methods)       │
│  • Generates examples with realistic sample values      │
│  • Creates point-wise "What It Does" / "When To Use"    │
│  • Truncates long docs, limits bullet points            │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│  Web Server (server/server.go)                          │
│                                                         │
│  • Embedded templates + static files (go:embed)         │
│  • Routes: / (index), /pkg/:name, /search, /api/search │
│  • Kubernetes-style UI with sidebar navigation          │
│  • Live search dropdown with instant results            │
└─────────────────────────────────────────────────────────┘
```

## What You Can Contribute

### 🐛 Bug Fixes
- Fix parsing issues for specific Go packages
- Fix UI rendering problems
- Fix incorrect examples being generated

### ✨ New Features
- **Better examples** — improve `sampleValue()` in generator.go for more types
- **More vanity imports** — add entries to `popularPackages` in main.go
- **Dark mode** — add CSS dark mode toggle
- **Export to Markdown** — export generated docs as .md files
- **Offline mode** — cache previously fetched packages
- **Go module version picker** — let users select specific versions

### 📖 Documentation
- Improve README
- Add doc comments to Go code
- Improve generated explanations for specific patterns

### 🎨 UI/UX
- Improve mobile responsiveness
- Add keyboard shortcuts
- Improve search ranking
- Add syntax highlighting to code examples

## File Guide

| File | Purpose | What to change |
|------|---------|----------------|
| `cmd/godoceasy/main.go` | CLI entry point | Add new CLI flags, popular packages |
| `internal/fetcher/github.go` | Fetches package source | Add new vanity import resolvers |
| `internal/parser/ast_parser.go` | Parses Go AST | Fix parsing edge cases |
| `internal/docs/generator.go` | Generates friendly docs | Improve examples, explanations |
| `internal/server/server.go` | HTTP server + search | Add new routes, improve search |
| `internal/server/templates/*.html` | HTML templates | UI changes |
| `internal/server/static/css/style.css` | Styles | Design improvements |
| `internal/server/static/js/search.js` | Live search | Search UX improvements |

## Code Style

- **Go standard formatting** — run `gofmt` before committing
- **No external dependencies** — we use only Go stdlib
- **Templates use `go:embed`** — static files are embedded in the binary
- **HTML uses `safeHTML`** — template function for rendering generated HTML

## How to Add a New Vanity Import

Edit `cmd/godoceasy/main.go` and add to the `popularPackages` map:

```go
var popularPackages = map[string]string{
    // ... existing entries ...
    "your.vanity/path": "https://github.com/org/repo.git",
}
```

## How to Improve Generated Examples

Edit `internal/docs/generator.go`:

1. **Better sample values** — update `sampleValue()` function
2. **New patterns** — update `isBuilderMethod()` to detect more patterns
3. **Better explanations** — update `generateFriendlyFunction()`

## Pull Request Process

1. Fork → branch → make changes
2. Run `go build ./cmd/godoceasy && go vet ./...`
3. Test with at least 3 different packages: `fmt`, `net/http`, and one remote package
4. Open a PR with a clear description of what changed and why

## Code of Conduct

Be kind. Be helpful. We're all here to make Go easier for beginners.

