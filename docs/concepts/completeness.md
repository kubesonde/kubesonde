# Probe completeness (quiescence)

Kubesonde probes **continuously** — it does not run once and stop. So "is it done?" is really "have the results stabilized enough to be worth reading?" Kubesonde answers this with **quiescence**.

## What "complete" means

Probing is considered complete when the recorded probe count has been **stable for a fixed window** (30 seconds by default). In other words: if no new probes have been recorded for the whole window, the results have quiesced and it's a good time to fetch them.

This is a *hint*, not a hard stop. The controller keeps probing after it reports `complete` — connectivity can change, and Kubesonde keeps observing.

## Where completeness shows up

**On the resource status** — as a `complete` flag and a `Complete` condition, so you can block on it:

```bash
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
```

```bash
kubectl get kubesondes
```

```
NAME                 NAMESPACE   PROBES   COMPLETE   AGE
kubesonde-sample     default     228      true       5m
```

**On the HTTP API** — query it directly:

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

- `window` is the quiescence window in **nanoseconds** — `30000000000` = 30 seconds.
- `secondsSinceLastChange` is how long the count has been unchanged. When it exceeds the window, `complete` becomes `true`.

## Why wait for it

Fetching results immediately after applying a `Kubesonde` resource may return empty or incomplete data, because probing has barely started. Always wait for the `Complete` condition (or a few minutes, depending on the number of pods) before you [fetch results](/how-to/fetch-results).

::: tip
The time to reach completeness scales with the number of pods and the number of probes. A large namespace takes longer; you can speed it up by [raising probe concurrency](/how-to/tune-concurrency).
:::
