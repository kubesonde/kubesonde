---
layout: home

hero:
  name: Kubesonde
  text: See the connectivity you actually have
  tagline: Probe live pods, record every connection attempt, and compare real network reachability against the policies you meant to enforce.
  actions:
    - theme: brand
      text: Quickstart
      link: /guide/quickstart
    - theme: alt
      text: What is Kubesonde?
      link: /guide/introduction
    - theme: alt
      text: View on GitHub
      link: https://github.com/kubesonde/kubesonde

features:
  - title: Runtime probing, not manifest reading
    details: Kubesonde instruments live pods and tries the connections, so you learn what is really reachable — not just what the NetworkPolicy YAML claims.
  - title: Intended vs. actual
    details: Surface the gap between the isolation you designed and the connectivity your CNI and policies actually enforce, before an attacker does.
  - title: Assert and audit
    details: Declare the outcomes you expect (Allow / Deny) and let Kubesonde verify them across every pod-to-pod pair in a namespace.
  - title: Visualize the graph
    details: Explore the resulting connectivity graph in the bundled UI or the hosted site — no manual reasoning about dozens of microservices.
---

<div class="ks-home-extra">

## Where to go next

- **[Tutorials](/tutorials/first-scan)** — learning-oriented, start-to-finish lessons. Begin here if Kubesonde is new to you.
- **[How-to guides](/how-to/install)** — task-oriented recipes for a specific goal ("target these pods", "export results").
- **[Concepts](/concepts/how-it-works)** — understanding-oriented background: how probing works and why.
- **[Reference](/reference/kubesonde-resource)** — information-oriented specs for the `Kubesonde` resource, HTTP API, and CLI.

::: tip New to Kubesonde?
Read the **[Introduction](/guide/introduction)** for the two-minute overview, then run the **[Quickstart](/guide/quickstart)** against a test cluster.
:::

</div>
