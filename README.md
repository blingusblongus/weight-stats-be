# weight-stats-be

Go backend for ingesting smart scale CSV exports into SQLite and serving them via REST API.

## How it works

- Watches a configured directory for `.csv` files (designed for syncthing)
- Parses scale data (weight, body fat, lean mass, BMR, etc.) and inserts into SQLite
- Deduplicates on timestamp via `INSERT OR IGNORE`
- Sends a single [ntfy](https://ntfy.sh) notification per batch with the count of new measurements
- Serves a simple REST API for the frontend

## API

| Method | Path | Params | Description |
|--------|------|--------|-------------|
| GET | `/api/measurements` | `start`, `end` (optional, YYYY-MM-DD) | All measurements in range, ordered by date |
| GET | `/health` | — | Health check |

## Configuration

Environment variables (or `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8083` | HTTP listen port |
| `DB_PATH` | `./weight-stats.db` | SQLite file path |
| `WATCH_DIR` | `./watch` | Directory to watch for CSVs |
| `NTFY_TOPIC` | — | ntfy endpoint (e.g. `https://ntfy.example.com/weight`) |
| `GIN_MODE` | `debug` | Set to `release` in production |

## CSV format

Expects exports from the smart scale app with headers:

```
Time,Weight,BMI,Body Fat,Fat-Free Body Weight,Subcutaneous Fat,Visceral Fat,Body Water,Muscle Mass,Skeletal Muscles,Bone Mass,Protein,BMR,Metabolic Age
```

Values include units (`228.9 lb`, `30.4 %`, `2134 kcal`) which are stripped during parsing.

## Deploy

Containerized via GHCR. See `docker-compose.yml` in the frontend repo for the full stack setup. The backend needs two volume mounts:

- `/data` — SQLite database (persistent)
- `/watch` — syncthing folder with CSV exports (read-only)

## Local dev

```bash
cp .env.example .env  # edit as needed
go run .
```
