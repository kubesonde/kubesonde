# Kubesonde Custom Resource Definition

Available commands: 

### `make test`
Runs unit tests

### `make run-test-env`
Runs a test environment for Kubesonde


### `make docker-build` and `make docker-push`
Push new versions of the artifacts to the upstream registry

### `make artifact`
Builds a new version of Kubesonde ready to be applied to a running k8s cluster
---
## About Kubernetes, docker containers and network namespaces 

In Kubernetes we can access the `containerID` property of a container. This maps the id of the docker container running the application.

> Q: Can I know if there are bogus containers running in docker? 
> 
> A: Not easy at the docker layer

## Issues
---
**Docker does not expose APIs to get the ID/location of the network namespace where containers run.**

**It is not obvious to get namespace information from docker containers because of *distroless containers* and limited capabilities that containers might have been assigned.**


**Typically, in cloud applications, developers do not have direct access to the underlying container runtime so they cannot tamper with it.**

In the host OS, I can check a single docker container and get its PID. From there, I can fetch the network namespace. This means that I know the network namespace of a given Pod. 

With little effort

```sh
netns=<netns>
sudo find -L /proc/[1-9]*/task/*/ns/net -samefile /run/netns/"$netns" | cut -d/ -f5
```

I can recover all the processes running on the same namespace. But this is not enough to know if there is a bogus container because, unfortunately, there might be muliple processes running on the same container. The developer may be able to see if there has been tampering with the original containers. 

> Partial solution: for each container, get namespace, run a shell and collect the network namespace information. 

Alternative: for each process of each container get the network namespace by accessing them (how easy with *distroless containers*)
