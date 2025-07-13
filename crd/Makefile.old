
SHELL := /bin/bash
# Image URL to use all building/pushing image targets
IMG ?= jackops93/kubesonde_controller:latest
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
	minikube delete && minikube start --force --memory=15000 --cpus=11 --driver=docker --kubernetes-version=v1.22.12 --feature-gates=EphemeralContainers=true --cni=calico
	minikube addons enable metrics-server
	echo "Installing prometheus..."
# cd kube-prometheus && kubectl apply -f manifests/setup > /dev/null && kubectl apply -f manifests > /dev/null
	make install
	make deploy-docker-container

.SILENT:
run-test-env-mac-os:
	echo "Setup MacOS test environment..."
	minikube delete
	minikube start  --memory=10000 --cpus=4 --kubernetes-version=v1.22.12 --driver=hyperkit --feature-gates=EphemeralContainers=true --cni=calico
	minikube addons enable metrics-server
# pushd kube-prometheus && kubectl apply -f manifests/setup  > /dev/null  && kubectl apply -f manifests  > /dev/null && popd
	make install
	make deploy-docker-container

.SILENT:
run-test-env:
	if [[ $$OS != *Darwin* ]]; then make run-test-env-mac-os; else make run-test-env-linux; fi

.SILENT:
deploy-docker-container:
	echo "Building kubesonde docker image..."
	eval $$(minikube -p minikube docker-env) && make docker-build && docker pull instrumentisto/nmap:latest &&  docker build ../docker/gonetstat -t monitor:latest
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

# Generate files to run kubesonde
artifact: manifests
	echo "Creating controller"
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default > kubesonde.yaml
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