FROM golang:1.15 AS builder

WORKDIR /tmp/src
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app .

FROM bash:5
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./

ENV HOSTNAME=.+
EXPOSE 8080

ENTRYPOINT ["bash"]
CMD ["-c", "/app stun -hp $HOSTNAME"]
