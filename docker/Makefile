IMG ?= jackops93/kubesonde_monitor:latest

all:
	docker build ./gonetstat -t ${IMG} > /dev/null
	docker push ${IMG} > /dev/null
