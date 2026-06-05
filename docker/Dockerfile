# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cloudstack ./cmd/cloudstack

# Final stage
FROM alpine:3.19

ARG VERSION=latest
ENV CLOUDSTACK_VERSION=${VERSION}
ENV CLOUDSTACK_STORAGE_PERSISTENT_PATH=/app/data

ENV AWS_DEFAULT_REGION=us-east-1
ENV AWS_ACCESS_KEY_ID=test
ENV AWS_SECRET_ACCESS_KEY=test
ENV AWS_CONFIG_FILE=/etc/cloudstack/aws/config

ENV GOSU_VERSION=1.17
ENV GOSU_AMD64_SHA256=bbc4136d03ab138b1ad66fa4fc051bafc6cc7ffae632b069a53657279a450de3
ENV GOSU_ARM64_SHA256=c3805a85d17f4454c23d7059bcb97e1ec1af272b90126e79ed002342de08389b

RUN set -eux; \
    apk add --no-cache shadow ca-certificates wget; \
    adduser -D -u 1001 cloudstack; \
    arch="$(apk --print-arch)"; \
    case "$arch" in \
        x86_64) gosuArch='amd64'; gosuSha256="$GOSU_AMD64_SHA256" ;; \
        aarch64) gosuArch='arm64'; gosuSha256="$GOSU_ARM64_SHA256" ;; \
        *) echo >&2 "unsupported arch: $arch"; exit 1 ;; \
    esac; \
    wget -q -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-${gosuArch}"; \
    echo "${gosuSha256}  /usr/local/bin/gosu" | sha256sum -c -; \
    chmod +x /usr/local/bin/gosu

WORKDIR /app

RUN mkdir -p /app/data /etc/cloudstack/aws \
    && printf '[default]\nendpoint_url = http://localhost:4566\n' > /etc/cloudstack/aws/config \
    && chown -R 1001:root /app \
    && chmod -R "g+rwX" /app

COPY --from=builder --chown=1001:root /app/cloudstack /app/cloudstack
COPY --chmod=755 docker/entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY --chmod=755 docker/localstack-parity.sh /usr/local/bin/localstack-parity.sh

VOLUME /app/data

EXPOSE 4566

HEALTHCHECK --interval=5s --timeout=3s --retries=5 \
    CMD wget -qO- http://localhost:4566/_cloudstack/health || exit 1

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["/app/cloudstack"]
