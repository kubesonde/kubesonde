IMG_MONITOR ?= jackops93/kubesonde_monitor:latest
IMG_PROBER ?= jackops93/kubesonde_prober:latest

all:
	docker build ./gonetstat -t ${IMG_MONITOR} > /dev/null
	docker build ./prober -t ${IMG_PROBER} > /dev/null
	docker push ${IMG_MONITOR} > /dev/null
	docker push ${IMG_PROBER} > /dev/null
