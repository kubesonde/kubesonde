# How Kubesonde works

Kubesonde determines connectivity **empirically**: instead of parsing `NetworkPolicy` manifests and reasoning about what *should* be allowed, it instruments live pods and actually attempts connections, recording what happens.

## The core idea

> Reading the manifests tells you what you declared. The only way to know what a pod can actually reach is to try.

Kubesonde tries. For the pods in a target namespace, it attempts connections between them (and to the ports you care about) and records, for each attempt, whether it succeeded or was refused. The collection of these attempts is a **connectivity graph** of the *actual* enforced state.

## The probing loop, at a high level

1. You create a [`Kubesonde` resource](/reference/kubesonde-resource) naming a target namespace and a probing mode (`all` or `none`, plus `include`/`exclude`).
2. The controller discovers the pods in that namespace and computes the set of connections to probe.
3. It dispatches probes through a bounded [worker pool](/how-to/tune-concurrency) (10 in parallel by default).
4. Each probe's outcome is recorded, along with the source, destination, port, and protocol.
5. The controller keeps probing continuously; a [`Complete` condition](/concepts/completeness) signals when results have stabilized enough to read.
6. You [fetch the results](/how-to/fetch-results) over the controller's HTTP API and [visualize](/how-to/visualize) them.

## Runtime, not static

Because probing happens at runtime against live pods, the results reflect everything that influences real traffic — the CNI plugin, `NetworkPolicies`, admission controllers, and how they interact — not just a single manifest. That's the whole point: interactions between these layers are where the gap between intended and actual connectivity hides.

## Assertions

If you attach an [`expected` outcome](/how-to/assert-outcomes) to a probe, Kubesonde compares the observed result to your expectation, turning a scan into a connectivity test.

## Related concepts

- [Intended vs. actual connectivity](/concepts/intended-vs-actual) — the problem Kubesonde exists to solve.
- [Architecture](/concepts/architecture) — the controller, CRD, and frontend.
- [Probe completeness](/concepts/completeness) — why the controller keeps probing.
