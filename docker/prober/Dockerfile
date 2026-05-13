FROM alpine:3.20

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.49

RUN apk add --no-cache \
    bash \
    curl \
    bind-tools \
    iputils \
    nmap \
    ca-certificates

RUN ARCH=$(uname -m) && \
    case "$ARCH" in \
        x86_64) ARCH=amd64 ;; \
        aarch64) ARCH=arm64 ;; \
        armv7l) ARCH=armv7 ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac && \
    curl -fsSL -o /usr/local/bin/grpc_health_probe \
      "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-${ARCH}" && \
    chmod +x /usr/local/bin/grpc_health_probe

CMD ["/bin/bash"]