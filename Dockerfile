FROM golang:1.22 AS builder

WORKDIR /tmp/src

COPY go.* .
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o /app

FROM bash:5
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

ENV HOSTNAME=.+
EXPOSE 8080

ENTRYPOINT ["bash"]
CMD ["-c", "/app stun -hp $HOSTNAME"]
