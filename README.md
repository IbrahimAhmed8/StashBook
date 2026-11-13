# StashBook

StashBook is an automated aggregator of **100% free games and tech books**
across storefronts — live at [stashbook.online](https://stashbook.online).

A Go engine scans several sources every 30 minutes, stores what it finds in
SQLite, and publishes static JSON that the React frontend renders. There is
no backend server and no user accounts: the whole site is a static build fed
by a scheduled job.

## Features

- **Live free offers** from Steam, Epic Games, GOG, Prime Gaming, itch.io and
  more, grouped by platform, with an **Expired** vault of past offers.
- **Iconic Drops**: deep discounts (50%+) on acclaimed titles.
- **Tech Archives**: free programming and developer books.
- **Game detail pages** with descriptions, end dates and outbound store links.
- **SEO-first**: per-game prerendered HTML, JSON-LD, sitemap and meta tags.

## Architecture

```
┌──────────────────────┐      ┌──────────┐      ┌────────────────────────┐
│  stashbook-engine    │      │ SQLite   │      │  public/               │
│  (Go, 8 clients) ────┼─────▶│ working  │─────▶│  games.json            │
│  Epic, Steam,        │      │ database │      │  books.json            │
│  CheapShark,         │      └──────────┘      │  sitemap.xml           │
│  GamerPower, dBooks  │                        └───────────┬────────────┘
└──────────────────────┘                                    │
                                                            ▼
                                              ┌────────────────────────┐
                                              │  React 19 + Vite     │
                                              │  (SPA, prerendered   │
                                              │   at build time)     │
                                              └────────────────────────┘
```

### Engine (`stashbook-engine/`)

One client per source, all running concurrently each run:

| # | Client | Source |
|---|--------|--------|
| 1 | Epic Games | Epic Store GraphQL (weekly free games) |
| 2 | Steam Live Deals | Steam search + app details (100% off) |
| 3 | CheapShark | CheapShark API ($0 deals) |
| 4 | Self-Tracker | Steam top sellers, price-history based F2P detection |
| 5 | Prime Gaming | GamerPower API (Amazon Prime giveaways) |
| 6 | Global Indie & Consoles | GamerPower API (GOG, itch.io, …) |
| 7 | Iconic Deals Hunter | Steam top sellers + CheapShark (50%+ off) |
| 8 | Tech Vault | dBooks API (free tech books) |

Each run marks all games expired, then re-activates the offers that are still
live, so disappearing offers land in the Expired vault. The layout inside the
engine:

- `main.go` — thin entry point: open the store, run the clients, write assets.
- `internal/sources/` — one file per source client.
- `internal/store/` — SQLite schema, reads and writes.
- `internal/assets/` — writes `games.json`, `books.json` and `sitemap.xml`.
- `internal/model/`, `internal/sanitize/`, `internal/slug/`, `internal/run/` —
  shared types and small helpers.

### Frontend (repo root)

React 19 + Vite + Tailwind CSS, styled retro (Press Start 2P). Pages are
lazy-loaded; `scripts/prerender.mjs` runs after `vite build` and generates
static HTML per route and per game for crawlers.

### Automation

- **Scraper** (`.github/workflows/scraper.yml`): runs the engine every 30
  minutes and commits updated `games.json`, `books.json` and `sitemap.xml`.
  The working database is **not committed** — it is restored from the previous
  run via `actions/cache`. A cold start (no cache) still works: the engine
  recreates the schema and repopulates current offers; only accumulated
  history (expired offers, tracked prices) is lost.
- **CI** (`.github/workflows/ci.yml`): on pull requests, runs `go vet` and
  `go test` for the engine plus `npm run lint`, `npm test` and `npm run build`
  for the frontend.
- **Deployment**: Cloudflare Pages is connected to GitHub; merging to `main`
  redeploys the live site.

## Repository layout

```
.
├── src/                    React frontend (pages, components, hooks, utils)
├── public/                 Static assets + generated games.json, books.json, sitemap.xml
├── scripts/                Build-time prerender script
├── stashbook-engine/       Go aggregation engine
│   ├── main.go
│   └── internal/           model, store, sources, assets, sanitize, slug, run
└── .github/workflows/      scraper.yml (scheduled job), ci.yml (pull requests)
```

## Run locally

### Frontend

```bash
npm install
npm run dev      # development server
npm run build    # production build + prerender
npm run lint
npm test
```

### Engine

```bash
cd stashbook-engine
go run main.go
```

This needs network access (it calls the storefront APIs) and takes a few
minutes. It writes `../public/games.json`, `../public/books.json`,
`../public/sitemap.xml` and keeps its working database in
`./stashbook.db` (git-ignored; never commit it).

```bash
go vet ./...
go test ./...
```

## Contributing notes

- Do not commit `*.db` files: the engine database is a cache, not data.
- The scraper commits generated assets to `main` on a schedule; avoid
  hand-editing `public/games.json`, `public/books.json` or `public/sitemap.xml`
  — they are overwritten every run.
- Merging to `main` redeploys the live site, so land changes through pull
  requests and let CI pass first.

## Disclaimer

StashBook is an independent tool, not affiliated with, endorsed by, or
connected to Epic Games, Valve (Steam), Amazon (Prime Gaming), GOG, itch.io,
or any other storefront mentioned on the site. All game titles, images and
logos belong to their respective copyright holders.
