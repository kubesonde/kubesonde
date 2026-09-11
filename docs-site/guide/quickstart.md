# Quickstart

This guide takes you from an empty cluster to a visualized connectivity graph in four steps. It should take about ten minutes.

## Prerequisites

- A running Kubernetes cluster (v1.11.3+) and `kubectl` configured to talk to it. A local [kind](https://kind.sigs.k8s.io/) or [minikube](https://minikube.sigs.k8s.io/) cluster is fine.
- Permission to create namespaces, CRDs, and deployments (cluster-admin during evaluation is simplest).

## 1. Install Kubesonde

Grab the latest `kubesonde.yaml` installer manifest from the [releases page](https://github.com/kubesonde/kubesonde/releases) and apply it:

```bash
kubectl apply -f kubesonde.yaml
```

This creates the `kubesonde-system` namespace, the Kubesonde CRDs, and the controller that runs the probes. It also deploys the frontend UI alongside the controller.

## 2. Create a Kubesonde resource

Once the controller is running, create a `Kubesonde` resource that describes what to probe. This example targets every pod in the `default` namespace:

```yaml
# probe.yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: kubesonde-sample
spec:
  namespace: default
  probe: all
```

Apply it:

```bash
kubectl apply -f probe.yaml
```

The controller starts probing the target namespace. Watch progress — including whether probing has quiesced — with:

```bash
kubectl get kubesondes
```

```
NAME                 NAMESPACE   PROBES   COMPLETE   AGE
kubesonde-sample     default     228      true       5m
```

## 3. Fetch the results

Wait for the `Complete` condition before fetching, so you don't read empty or partial data:

```bash
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
```

Then port-forward the controller and download the recorded probes:

```bash
kubectl --namespace kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
curl localhost:2709/probes > probes.json
```

::: warning
Fetching results immediately after applying the resource may return empty or incomplete data. Wait for the `Complete` condition (or a few minutes, depending on the number of pods). See [Probe completeness](/concepts/completeness) for why.
:::

## 4. Visualize

The installer deploys the Kubesonde frontend inside the cluster. Port-forward the frontend service and the controller (the UI reads from `localhost:2709`), then open the UI:

```bash
kubectl --namespace kubesonde-system \
  port-forward service/kubesonde-frontend 8088:8088 &
kubectl --namespace kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709:2709 &
```

Browse to **[http://localhost:8088](http://localhost:8088)**.

Alternatively, upload the `probes.json` you exported to the [hosted Kubesonde website](https://kubesonde.jackops.dev) to explore the connectivity graph.

## Clean up

```bash
kubectl delete -f probe.yaml     # remove the scanner
kubectl delete -f kubesonde.yaml # remove the controller and CRDs
```

