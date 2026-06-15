# Engineering Notes

Interesting decisions, deviations, and trade-offs behind this submission.

---

## Ingest

**Idempotency / dedup key = `(client_id, measured_at)`.** A unique composite index enforces "one reading per station per timestamp," and the insert uses `ON CONFLICT DO NOTHING`, so a sensor retrying a delivery is a no-op — not a duplicate row, not an error. Chosen over SELECT-then-INSERT (races under concurrent delivery) and over a plain insert (errors/duplicates on redelivery).

**201 vs 200 makes dedup observable.** A fresh insert returns 201; a redelivered duplicate (skipped by `ON CONFLICT`) returns 200. The controller reads `RowsAffected` (1 = inserted, 0 = skipped) into an internal `Created` flag (`json:"-"`). Both are success — a 409 would provoke retry storms from flaky firmware — but the split lets an operator see deduplication on the wire.

**Validation: 400 vs 422, and missing ≠ 0.** `pm2_5`/`pm10` must be finite, `>= 0`, `<= ~1000` (a sanity ceiling guarding the KPIs from sensor garbage). Out-of-bounds / missing fail struct-tag validation → **400** (the shared `router.NewValidator`); non-finite (`NaN`/`Inf`) and malformed JSON fail binding → **422**. The fields are pointers, so an absent PM is distinguishable from a real `0` and rejected rather than silently accepted. I kept the shared validator's 400 rather than forcing 422 on this one route, to match how enrollment already validates.

**Unknown `client_id` → 404.** Ingest never auto-creates stations; enrollment stays the single path that creates them.

**Tested against a real DB, not mocks.** The ingest path is covered against an in-memory SQLite with true migrations, the dedup index, and the actual `ON CONFLICT` upsert, and driven through the full chi router via `httptest` — so the 201/200 and 400/422/404 status contract is asserted end to end.

---

## Data model

**Lean readings row (deviation from the skeleton's `Model` base).** `Reading` carries only `client_id, pm2_5, pm10, measured_at, created_at` — deliberately no `updated_at`, no soft-delete `deleted_at`. Readings are immutable append-only facts; soft-delete would add an indexed column and a `WHERE deleted_at IS NULL` to every read for no benefit at fleet scale. A conscious departure from the enrollment model's convention.

**One index does two jobs.** A single unique composite index on `(client_id ASC, measured_at DESC)` both enforces the dedup constraint (its uniqueness backs the `ON CONFLICT` upsert) and orders the per-station trend query's keyset pagination plus the 24h connection scan — no separate sort. A second ASC index would only have duplicated the column pair and added write-path cost.

**SQLite concurrency tuning (measured).** Under concurrent ingest the default unlimited pool had connections racing SQLite's single write lock; some timed out with `database is locked`. The DSN now sets `journal_mode=WAL`, `busy_timeout=5000`, `txlock=immediate`, `synchronous=NORMAL`, with the pool capped at 4. A 10k-station simulator run went from 4 enroll 500s + 859ms p95 to **0 failures, ~50ms p95**, throughput ~7×. Postgres/MySQL would raise the pool.

---

## Display / query API

**Latest-per-station via a denormalized cache (measured).** Each station's most recent reading lives in a `latest_readings` table (one row per station), updated on ingest via a guarded upsert (`ON CONFLICT(client_id) DO UPDATE ... WHERE measured_at < <new>`) so a late or duplicate delivery is a no-op. `/v1/stations` and `/v1/stats` read it with an O(stations) join instead of a window scan over the whole readings table (O(total), ~2.9M rows/day at 10k stations). The cache write shares the ingest transaction, so the two can't diverge. A 10k-station run confirmed **0 mismatches** vs a full recompute, and the snapshot stays flat as history grows (the window query degraded to 459ms at 49k rows and climbing). `/v1/clients/{id}/readings` stays raw cursor-paginated rows with a `since`/`from`+`to` filter — no server-side bucketing (see *With more time*).

**KPI definitions are fleet policy, not hardcoded logic.** *Active Sensors* = stations with a reading in the last N min; *Avg PM2.5* = mean latest reading across active stations; *Poor Connection* = count with received/expected ratio < threshold; *Network Stability* = fleet-mean ratio. Windows, expected cadence, and threshold are `display.*` config (default: expected = 1440 ÷ 5 = 288/day, threshold 0.8, active window 15 min). Active/Avg use the short "is it reporting *now*?" window; Poor/Stability use the 24h "has it reported *reliably*?" window.

> **Why the demo shows everything "Poor" / single-digit stability.** Expected is calibrated for 24h (288 readings). The simulator streams a short burst (~20–30 readings), so each ratio is `received ÷ 288` ≈ a few percent — correctly below 0.8. Stream ~5 min (≈1 reading/sec passes 288) or lower `connection_window_minutes` to see healthy values. The metric is behaving as defined, not a bug.

**Districts derived backend-side from coordinates.** A station's district comes from point-in-polygon of its `(lat, lng)` against a static Berlin Bezirke GeoJSON loaded once at startup — the frontend does no geo math.

- *Computed once, cached.* Coordinates never change, so district is a stable fact held in an in-memory `clientID → district` map, not recomputed per request (a per-request sweep is ~1.5×10⁸ comparisons at 10k stations). Each lookup is bounding-box pre-filtered — the full-border test runs only if the point is inside a district's bounding rectangle.
- *Point-in-polygon only, via `paulmach/orb`.* `orb`'s `geojson` sub-package pulls in a MongoDB BSON driver, so the file is decoded with stdlib `encoding/json` into orb geometry instead — keeping the vendored footprint stdlib-only.
- *Sensors outside Berlin show "Unknown" — by design.* A sensor's district is found by matching its coordinates against Berlin's 12 real district shapes. The simulator drops sensors into a rectangle around the city, but Berlin isn't a rectangle — so about 40% land just outside it, in the neighbouring Brandenburg countryside. Those fall in no district and correctly come back as "Unknown" rather than being given a wrong guess. (So a high share of "Unknown" under the simulator is the lookup working, not a bug.)

---

## Frontend

**Styling vs logic split.** Style *values* (colors, spacing, radii) are extracted once from the mockup into a tokens file + a PrimeVue theme preset; component *logic* (state, data flow, the value→band mapping, reactivity) is hand-written. Components carry no inline color or padding literals.

**Lite states are correctness, not polish.** On a live ingest pipeline a blank view is ambiguous — loading vs empty vs backend-down look identical. Every fetch branches skeleton / empty / error (map, watchlist, KPIs, chart).

**Per-station readings = composable, not store.** The chart's readings are view-local and parameterized (station + range), so they live in a `useStationReadings` composable, not the global store — which holds only genuinely shared state (the fleet, stats, the selected station). Selecting a marker sets the shared selection; the chart watches it and refetches.

---

## With more time

Each extends an existing piece rather than adding a new domain.

- **Server-side range aggregation** — hourly/daily bucketing on the readings endpoint. The 24H/7D/30D selector ships, but the 50-row newest-first page cap means a densely-reporting station renders the same recent ~50 points for 7D/30D as for 24H; bucketing is what makes long ranges show a genuine multi-day shape.
- **Batched watchlist sparkline feed** — the sparkline is now lightweight inline SVG (down from ~50 Chart.js canvases), but each row still fetches its own 24h readings (N requests). At fleet scale this needs one batched/summary endpoint.
- **Live updates** — the dashboard fetches once on mount, so the fleet snapshot (markers, KPIs, watchlist) is frozen until a manual refresh; a short polling interval — or SSE for true push — would keep it live. (The trend chart already refetches on selection/range change.)
- **Map marker clustering** — keep ~10k markers readable and responsive.
- **Search debounce / typeahead** — basic district/sensor-ID search ships; the debounced refinement is deferred.
- **Persist the derived district** — it's cached in memory today; writing it on the row would survive restarts without the startup sweep.
- **Watchlist auto-scroll to selected row** — selection pans the map (shipped); scrolling the table to the matching row needs DOM `scrollIntoView` against PrimeVue internals.
