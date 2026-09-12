# Intended vs. actual connectivity

Kubernetes network security is **declarative**. You write `NetworkPolicies` describing the connectivity you intend, and you trust the cluster to enforce it. Kubesonde exists because the *intended* state and the *actual* enforced state are not always the same — and the difference is where security problems live.

## The gap

| Intended connectivity | Actual connectivity |
| --- | --- |
| What your `NetworkPolicy` manifests describe | What pods can really reach right now |
| Derived by reading YAML | Derived by attempting connections |
| Assumes correct, complete policies | Reflects CNI behavior, policy interactions, and mistakes |

Three recurring reasons the gap opens:

1. **NetworkPolicies are hard to get right.** A missing selector, an overly broad `podSelector`, or a forgotten default-deny policy can silently leave pods reachable. Reading the manifests does not reveal this.
2. **Enforced state can differ from declared state.** CNI plugins, admission controllers, and overlapping policies interact in non-obvious ways.
3. **Manual auditing does not scale.** In a namespace with dozens of microservices, reasoning by hand about which pod can talk to which — and on which ports — is error-prone and quickly impossible.

## How Kubesonde closes the loop

Kubesonde measures the actual state directly by [probing live pods](/concepts/how-it-works). You can then:

- Compare the observed graph against your mental model or your policy manifests.
- [Assert expected outcomes](/how-to/assert-outcomes) so the comparison is automatic and repeatable.
- Catch unintended reachability **before** an attacker finds it.

The goal is not to replace your policy engine — it is to give you evidence about whether that engine is doing what you think.
