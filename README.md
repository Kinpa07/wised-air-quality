# Wised Berlin Air-Quality Monitor

A Go sensor-readings service (`main-service/`) and a Vue 3 dashboard
(`frontend/`) for monitoring Berlin's air-quality sensor fleet. Sensors POST
PM2.5/PM10 readings to the service; the dashboard shows a live Berlin map, a
districts watchlist, per-station trend charts, and fleet KPIs.

See [`NOTES.md`](NOTES.md) for design decisions and trade-offs.

## Prerequisites

- **Docker** — runs the backend and the provided fleet simulator (no local
  Go/C toolchain required).
- **Node.js 18+ and npm** — runs the dashboard.
- _(Optional)_ **Go 1.26 + a C compiler** — only to run the backend without
  Docker (SQLite needs CGO).

## Quick start

Three terminals.

### 1. Backend — `:8080`

```bash
cd main-service
docker build -t wised .
docker run -p 8080:8080 wised
```

> No Docker? With Go 1.26 + a C compiler: `cd main-service && go run .`

### 2. Simulator — streams readings into the backend

The fleet simulator is the Docker image provided with the assignment
(`wised-interview-simulator.tar`); it isn't included in this repo. Place it in
`main-service/`, then load and run it:

```bash
docker load -i main-service/wised-interview-simulator.tar
docker run --rm wised-interview-simulator:latest
```

Enrolls stations and streams readings at `http://host.docker.internal:8080`.
Leave it running ~5 minutes so the connection KPIs (Poor Connection, Network
Stability) reach healthy values — they're calibrated for a full reporting
window. _(On Linux, run with `--network=host`.)_

Many sensors will show district **"Unknown"** — the simulator scatters them
across a rectangle around Berlin, so ~40% land outside the city's real
districts. That's expected, not a bug (see [`NOTES.md`](NOTES.md)).

### 3. Dashboard — `:5173`

```bash
cd frontend
npm install
npm run dev
```

Open <http://localhost:5173>. Vite proxies `/api` to the backend on `:8080`.

## Tests

```bash
cd main-service
go test ./...
```
