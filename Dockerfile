FROM golang:latest AS builder
WORKDIR /app
COPY srcs/ .
RUN go mod tidy && go build -o inquisitor .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libpcap-dev libc6 && apt-get clean
WORKDIR /app
COPY --from=builder /app/inquisitor /app/inquisitor
RUN chmod +x inquisitor
CMD ["/bin/bash"]