FROM golang:1.13 AS builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/kubatek94/dyndns-cloudflare-adapter
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./

EXPOSE 8080
ENTRYPOINT ["./app"]

