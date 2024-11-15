![workflow](https://github.com/kubesonde/kubesonde/actions/workflows/go_main.yaml/badge.svg)
![frontend_test](https://github.com/kubesonde/kubesonde/actions/workflows/frontend_dev.yaml/badge.svg)
![frontend_deployment](https://github.com/kubesonde/kubesonde/actions/workflows/deploy_frontend.yaml/badge.svg)
[![Netlify Status](https://api.netlify.com/api/v1/badges/df3643ab-e317-4b96-b5c2-de937837b375/deploy-status)](https://app.netlify.com/sites/testksonde/deploys)



![Kubesonde logo](frontend/public/logo257.png "Kubesonde logo")

# Kubesonde

Kubesonde is a tool to probe and test network security policies in Kubernetes.

![kubesonde infra](docs/kubesonde.png "kubesonde infrastructure")

## Structure of the project
Folders are organized as follows: 
- `crd`: backend service and kubesonde CRD 
- `docs`: documentation of the project/ideas.
- `frontend`: contains the UI for analyzing the probe outputs
- `examples`: sample output from Kubesonde

## Run Kubesonde
### 1. Start the Kubernetes engine

You can run Kubernetes on the cloud, bare-metal or via Minikube or Kind.
### 2. Install the app to test

Install the application you want to test (e.g., `helm install wordpress bitnami/wordpress`). Make sure that the app is running with no errors.

### 3. Install Kubesonde

To install kubesonde run `kubectl apply -f kubesonde.yaml`. This creates all the required resources to run Kubesonde on your cluster. After that, you can install a scanner object for Kubesonde. The following is a Kubesonde object example that targets the default namespace: 
```yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: kubesonde-sample
spec:
  namespace: default
  probe: all
```
You can save it in a file `probe.yaml` and then apply it with `kubectl apply -f probe.yaml`
### 4. Fetching the results

To fetch the results, you need to use the following commands:

`kubectl --namespace kubesonde port-forward deployment.apps/kubesonde-controller-manager 2709`. This command creates a port mapping between your local computer and the Kubesonde deployment.

`curl localhost:2709/probes > <output-file>.json`. This command gets the probe result and stores it in an output file.

:warning: If you try to get the results of the probe just after applying it in the cluster the results may be empty or incomplete. Wait a few minutes (depending on the amount of pods) to get better results.
### 5. View results

Navigate to the [kubesonde website](https://kubesonde.jackops.dev) and upload the generated file to see the results.
 
## Credits

Logo from [Elisabetta Russo](stelladigitale.it) info@stelladigitale.it