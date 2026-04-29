FROM golang:1.22-bookworm AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends libpcap-dev && rm -rf /var/lib/apt/lists/*
COPY srcs/ .
RUN go mod tidy && go build -o inquisitor .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
        libpcap0.8 iproute2 iputils-ping ftp lftp tcpdump ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/inquisitor /usr/local/bin/inquisitor
RUN chmod +x /usr/local/bin/inquisitor
CMD ["/bin/bash"]
