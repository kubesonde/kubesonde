# Controller HTTP API

The controller (`kubesonde-controller-manager`) serves an HTTP API on port **`2709`**. Port-forward it to reach the endpoints locally:

```bash
kubectl --namespace kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
```

All endpoints send permissive CORS headers, so the browser-based [frontend](/how-to/visualize) can read them directly.

## `GET /probes`

Returns the recorded probe state as a JSON **object**. This is the payload you export and load into the UI.

```bash
curl -s localhost:2709/probes > probes.json
```

The response is an object, not a bare array. The probes live under `items`:

```json
{
  "items": [ { "type": "Probe", "resultingAction": "Allow", ... } ],
  "errors": [ { "reason": "command terminated with exit code 1", "value": { ... } } ],
  "podNetworking": [],
  "podNetworkingv2": { },
  "podConfigurationNetworking": { }
}
```

| Top-level field | Type | Description |
| --- | --- | --- |
| `items` | array | The recorded probes (see the fields below). |
| `errors` | array | Probes that failed to run, each with a `reason` and the offending probe under `value`. |
| `podNetworking`, `podNetworkingv2`, `podConfigurationNetworking` | — | Internal networking snapshots used by the UI. |

Each entry in `items` includes (fields are omitted when empty):

| Field | Type | Description |
| --- | --- | --- |
| `type` | string | `Probe` (a connection attempt) or `Information`. |
| `expectedAction` | string | The asserted outcome, if you set `expected` — `Allow` or `Deny`. |
| `resultingAction` | string | The observed outcome — `Allow` or `Deny`. |
| `source` | object | Origin endpoint (see below). |
| `destination` | object | Destination endpoint (see below). |
| `destinationHostnames` | string[] | Hostnames resolved for the destination. |
| `protocol` | string | Transport protocol. |
| `port` | string | Destination port (defaults to `80`). |
| `forwardedPort` | string | Local forwarded port, when applicable. |
| `timestamp` | int64 | Unix timestamp of the probe. |
| `debugOutput` | string | Raw output of the request (e.g. the HTTP code for TCP). |

Each endpoint (`source` / `destination`) may include `type` (`Pod`, `Service`, or `Internet`), `name`, `IPAddress`, `namespace`, `labels`, `deploymentName`, and `replicaSetName`.

## `GET /probes/status`

Reports whether probing has [quiesced](/concepts/completeness).

```bash
curl -s localhost:2709/probes/status
```

```json
{
  "complete": true,
  "count": 228,
  "window": 30000000000,
  "secondsSinceLastChange": 42.5
}
```

| Field | Type | Description |
| --- | --- | --- |
| `complete` | bool | `true` when at least one probe has been recorded and the count has been stable for `window`. |
| `count` | int | Number of recorded probe items. |
| `window` | int | Quiescence window as a duration in **nanoseconds** (`30000000000` = 30s). |
| `secondsSinceLastChange` | float | How long the count has been stable, or `-1` if no probes have been recorded yet. |

## `POST /probes/clear`

Clears the recorded probe state on the controller. Useful when you want to re-run a scan from a clean slate without restarting the controller.

```bash
curl -s -X POST localhost:2709/probes/clear
```

## Notes

- The controller keeps probing after `complete` becomes `true`; polling `/probes` again later can return more data.
- Only `GET` is accepted on the two read endpoints; other methods return `405 Method Not Allowed`.
