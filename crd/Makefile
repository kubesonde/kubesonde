
SHELL := /bin/bash
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

OS=$(uname -s)

all: manager

.SILENT:
### Testing recipes ###
# Run test environment
run-test-env-linux:
	echo "Setup linux test environment"
	minikube delete && minikube start --force --memory=15000 --cpus=5 --driver=docker --kubernetes-version=v1.22.12 --feature-gates=EphemeralContainers=true --cni=calico 
	echo "Installing prometheus..."
# cd kube-prometheus && kubectl apply -f manifests/setup > /dev/null && kubectl apply -f manifests > /dev/null
	make install
	make deploy-docker-container

.SILENT:
run-test-env-mac-os:
	echo "Setup MacOS test environment..."
	minikube delete
	minikube start  --memory=10000 --cpus=4 --kubernetes-version=v1.22.12 --driver=hyperkit --feature-gates=EphemeralContainers=true --cni=calico
# pushd kube-prometheus && kubectl apply -f manifests/setup  > /dev/null  && kubectl apply -f manifests  > /dev/null && popd
	make install
	make deploy-docker-container

.SILENT:
run-test-env:
	if [[ $$OS != *Darwin* ]]; then make run-test-env-mac-os; else make run-test-env-linux; fi

.SILENT:
deploy-docker-container:
	echo "Building kubesonde docker image..."
	eval $$(minikube -p minikube docker-env) && make docker-build && docker pull instrumentisto/nmap:latest &&  docker build docker/pynetstat -t netinfo:latest
	make deploy

test: generate fmt vet manifests
	go test ./... -v -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go  > /dev/null

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests 
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	echo "Installing kubesonde CRD in the cluster..."
	kustomize build config/crd | kubectl apply -f - > /dev/null

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	echo "Deploying controller"
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f - > /dev/null

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	echo "Generating controller..."
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases  > /dev/null

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."  > /dev/null

# Build the docker image
docker-build:
	docker build . -t ${IMG} > /dev/null

# Push the docker image
docker-push:
	docker push ${IMG} > /dev/null

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# test utils
wait-containers-and-launch-prometheus:
	kubectl wait --for=condition=Ready --timeout=10m pods --all --all-namespaces
	kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090 &

build-graph:
	echo "Waiting for kubesonde to finish probing..."
	sleep 180
	pushd scripts && pipenv run python produce_graph.py $(testcase) && popd
# tests

istio-test:
	kubectl create namespace istio-system
	helm install istio-base istio/base -n istio-system
	helm install istiod istio/istiod -n istio-system 
	kubectl create namespace istio-ingress
	kubectl label namespace istio-ingress istio-injection=enabled
	helm install istio-ingress istio/gateway -n istio-ingress
	kubectl label namespace default istio-injection=enabled
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/networking/bookinfo-gateway.yaml
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/networking/destination-rule-all.yaml


microservices-test: run-test-env deploy-docker-container
	if [[ ! -d '/tmp/microservices-demo' ]]; then git clone https://github.com/GoogleCloudPlatform/microservices-demo /tmp/microservices-demo; fi
	/tmp/istio-1.7.2/bin/istioctl manifest install --set profile=demo -y
	kubectl label namespace default istio-injection=enabled
	pushd /tmp/microservices-demo && kubectl apply -f ./release/kubernetes-manifests.yaml && popd
	make wait-containers-and-launch-prometheus
	kubectl apply -f config/samples/security_v1_kubesonde.yaml
	make build-graph testcase=microservices
	kubectl apply -f  test-without-internet-connection/disallow-internet.yaml
	make build-graph testcase=microservicesnointernet

pitstop-test: run-test-env deploy-docker-container
	if [[ ! -d '/tmp/pitstop' ]]; then git clone https://github.com/EdwinVW/pitstop /tmp/pitstop; fi
	pushd /tmp/pitstop/src/k8s && chmod +x start-all.sh && ./start-all.sh && popd
	make wait-containers-and-launch-prometheus
	kubectl apply -f config/samples/sample-pitstop.yaml
	make build-graph testcase=pitstop
	kubectl apply -f  test-without-internet-connection/disallow-internet.yaml
	make build-graph testcase=pitstopnointernet

waveworks-microservices-test: run-test-env deploy-docker-container
	if [[ ! -d '/tmp/waveworks-microservices' ]]; then git clone https://github.com/microservices-demo/microservices-demo /tmp/waveworks-microservices ; fi
	pushd /tmp/waveworks-microservices &&  kubectl create -f deploy/kubernetes/complete-demo.yaml && popd
	make wait-containers-and-launch-prometheus
	kubectl apply -f config/samples/sample-sock-shop.yaml
	make build-graph testcase=sockshop
	kubectl apply -f  test-without-internet-connection/disallow-internet.yaml
	make build-graph testcase=sockshopnointernet

guestbook-test: run-test-env deploy-docker-container
	if [[ ! -d '/tmp/examples' ]]; then echo "Installing Guestbook"; git clone https://github.com/kubernetes/examples.git /tmp/examples; fi
	echo "Guestbook installed"
	pushd /tmp/examples/guestbook &&  kubectl apply -f all-in-one/ && popd
	make wait-containers-and-launch-prometheus
	kubectl apply -f config/samples/sample-guestbook.yaml
	make build-graph testcase=examples
	kubectl apply -f  test-without-internet-connection/disallow-internet.yaml
	make build-graph testcase=examplesnointernet

complete-test: guestbook-test waveworks-microservices-test pitstop-test microservices-test istio-test
	echo "Done!"
