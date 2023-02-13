IMG ?= registry.cs.aalto.fi/kubesonde/monitor:latest

all:
	docker build ./gonetstat -t ${IMG} > /dev/null
	docker push ${IMG} > /dev/null
