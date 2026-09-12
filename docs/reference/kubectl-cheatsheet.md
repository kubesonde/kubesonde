# kubectl cheat sheet

A quick reference of the commands used throughout the Kubesonde docs.

## Install / uninstall

```bash
kubectl apply -f kubesonde.yaml          # install controller, frontend, CRDs
kubectl -n kubesonde-system get pods     # check the controller is running
kubectl delete -f kubesonde.yaml         # uninstall
```

## Manage a probe

```bash
kubectl apply -f probe.yaml              # create a Kubesonde resource
kubectl get kubesondes                   # list probes with count + complete
kubectl describe kubesonde kubesonde-sample
kubectl delete -f probe.yaml             # stop probing
```

## Wait for completeness

```bash
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
```

## Fetch results

```bash
kubectl -n kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
curl -s localhost:2709/probes > probes.json
curl -s localhost:2709/probes/status
curl -s -X POST localhost:2709/probes/clear   # reset recorded probes
```

## Visualize

```bash
kubectl -n kubesonde-system port-forward service/kubesonde-frontend 8088:8088 &
kubectl -n kubesonde-system port-forward deployment.apps/kubesonde-controller-manager 2709:2709 &
# open http://localhost:8088  — or upload probes.json to https://kubesonde.jackops.dev
```

## Tune concurrency

```bash
kubectl -n kubesonde-system set env \
  deployment/kubesonde-controller-manager KUBESONDE_PROBE_WORKERS=20
```

## Inspect the CRD

```bash
kubectl get crd kubesondes.security.kubesonde.io
kubectl explain kubesonde.spec
```
