# Feature Plan: API-driven live graph in the frontend

## Goal

Allow the frontend to load probe data directly from a running Kubesonde
controller instead of a manually uploaded JSON file. When the
`VITE_API_SERVER` environment variable is set, the UI automatically fetches
graph data from the controller's REST API, polls every ~10 seconds until the
probe run is complete, re-renders as new data arrives, and offers a manual
refresh button once the run is ready.

## Background: how it works today

- Data enters the UI **only via file upload** (`HomeComponent` +
  `GraphJSONUploadComponent`) using `use-file-picker`.
- Upload parses JSON into a `ProbeOutput`, then navigates to
  `/graph/:name` passing the parsed data through React Router `state`.
- `GraphFromLocation` reads `location.state.data`, builds nodes/edges, and
  renders `GraphBase` plus the stats / error / netinfo tables.
- There is **no HTTP/fetch layer** anywhere in the frontend today.
- Vite is the bundler, so env vars are exposed as `import.meta.env.VITE_*`.

## Controller endpoints (already exist, port 2709)

Defined in `crd/rest_apis/`:

- `GET /probes` — full `ProbeOutput` JSON (same shape the UI renders from a file).
- `GET /probes/status` — completeness signal:
  `{ complete: bool, count: int, window, secondsSinceLastChange }`.
  `complete` flips to `true` once the recorded probe count has been stable for a
  quiescence window. **This is the authoritative "ready" signal** — no
  frontend-side heuristic is required.
- `POST /probes/clear` — resets probe state (candidate for a future "re-run"
  button; out of scope here).

## Decisions locked in

- **Env var name:** `VITE_API_SERVER` (Vite only exposes `VITE_`-prefixed vars
  to the client).
- **Ready detection:** use `GET /probes/status` → `complete`. No guessing.
- **"Stop polling when complete" is frontend-only.** The controller keeps
  probing continuously in the background; the frontend simply stops its 10s
  timer once `complete` is true and treats the result as a snapshot. The manual
  refresh button is how the user pulls a newer snapshot afterwards.
- **Manual refresh button:** shown only *after* `/probes/status` reports ready.
- **Controller unreachable / not ready:** no elaborate error UX — just render a
  simple message ("Probes are not ready or not available") and let the poll /
  refresh retry.
- **Container + presentational split:** a container that fetches the data (via
  the polling hook / state) feeds a presentational `GraphView` that takes the
  data. This mirrors how `GraphFromLocation` works today, except the data comes
  from the server instead of router `state`.
- **Routing / landing:** when `VITE_API_SERVER` is set, the app **launches into
  the new live view** as its landing page. The sidebar stays on the left and
  still offers the existing "upload a file" home and the example probe. When the
  env var is absent, behaviour is unchanged.
- **CORS:** add permissive CORS headers in the Go handlers
  (`Access-Control-Allow-Origin: *`, `Allow-Methods: GET, POST, OPTIONS`, handle
  `OPTIONS`). No auth, so `*` is fine.
- **Tests for everything:** every new unit (hook, container, `GraphView`,
  config, CORS handlers) ships with tests.

## Implementation plan

### Frontend

1. **`frontend/src/utils/config.ts`**
   - Export `apiServer = import.meta.env.VITE_API_SERVER` and
     `isApiMode = !!apiServer`.
   - Add the env var's type declaration to `vite-env.d.ts`
     (or `react-app-env.d.ts`).

2. **`frontend/src/utils/useProbeData.ts`** (new polling hook)
   - Every 10s: `GET ${apiServer}/probes/status`.
   - While `complete === false`: keep polling; optionally fetch `/probes` to
     show progressive results, or surface a "probing… N probes so far" state
     using `count`.
   - Once `complete === true`: do a final `GET /probes`, **stop the interval**,
     expose the data.
   - Returns `{ data, status, ready, loading, error, refresh }`.
   - `refresh()` re-fetches on demand and restarts polling if not yet complete.
   - Clear the interval on unmount.

3. **Extract presentational `GraphView({ data, title })`**
   - Pull the rendering body out of `GraphFromLocation` into a reusable
     presentational component that takes the probe data as a prop.
   - **Important:** move the `cleanupProbeOutput(...)` call and the
     `netinfo2Table` transforms into `GraphView` (or apply them before passing
     data in) so the upload path and the API path produce identical output.
   - `GraphFromLocation` becomes a thin wrapper: read `location.state`, render
     `<GraphView data title />`.

4. **API container (`GraphFromApi`)**
   - A container component that fetches its own data via `useProbeData` (state,
     not props) and renders `<GraphView data title="Live cluster" />`.
   - While not ready: render the polling/progress indicator.
   - On fetch failure / not-ready: render the simple message
     "Probes are not ready or not available".

5. **Routing / landing**
   - When `isApiMode`, the `/` route renders `GraphFromApi` (the new live view is
     the landing page). Keep the sidebar; add links to the existing upload home
     and the example probe so both remain reachable.
   - When `!isApiMode`, `/` renders the current `HomeComponent` and routing is
     unchanged.

6. **Manual refresh button**
   - Appears once `/probes/status` reports ready (`ready === true`), on the
     live view, with a last-updated timestamp.
   - Clicking calls `refresh()` to re-fetch `/probes`.

### Backend (`crd/rest_apis/`)

7. **CORS headers**
   - Add `Access-Control-Allow-Origin: *`, `Access-Control-Allow-Methods:
     GET, POST, OPTIONS`, and explicit `OPTIONS` preflight handling to the
     `/probes`, `/probes/status`, and `/probes/clear` handlers, since the
     browser fetch to `:2709` is cross-origin. Ideally via one small shared
     middleware wrapping the handlers rather than duplicating headers.

## Task breakdown

Each task ships with its tests.

1. **Config + env var** — `config.ts` (`apiServer`, `isApiMode`), `vite-env.d.ts`
   typing. *Tests:* set/unset `VITE_API_SERVER`, trailing-slash normalization.
2. **`useProbeData` hook** — poll `/probes/status`, fetch `/probes`, stop on
   `complete`, `refresh()`, cleanup/abort on unmount. *Tests:* mocked fetch —
   polls until complete then stops, `refresh` re-fetches, error → error state,
   unmount aborts.
3. **`GraphView` extraction** — presentational component incl. `cleanupProbeOutput`
   + `netinfo2Table`. *Tests:* renders nodes/edges/tables from a sample
   `ProbeOutput`; `GraphFromLocation` still works via `GraphView`.
4. **`GraphFromApi` container** — wires hook → `GraphView`, progress + not-ready
   message. *Tests:* loading → ready → graph; unreachable → message.
5. **Routing / landing** — `/` renders `GraphFromApi` in API mode, `HomeComponent`
   otherwise; sidebar keeps upload + example links. *Tests:* landing route per
   mode.
6. **Refresh button** — appears after ready, calls `refresh()`, shows
   last-updated. *Tests:* hidden until ready; click triggers re-fetch.
7. **CORS middleware (Go)** — shared wrapper on the three handlers. *Tests:* Go
   handler tests assert CORS headers + `OPTIONS` preflight returns 200/204.

## Notes

- No new frontend dependencies required (native `fetch`).
- Suggested workflow: implement on a feature branch off `dev`.
- Local usage: `kubectl port-forward … 2709`, then set
  `VITE_API_SERVER=http://localhost:2709` for the dev server.
