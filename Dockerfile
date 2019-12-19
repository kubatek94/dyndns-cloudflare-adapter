FROM golang:1.13 AS builder

# # Download and install the latest release of dep
# RUN wget -O /usr/bin/dep https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 && \
#     chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/kubatek94/dyndns-cloudflare-adapter
# COPY Gopkg.toml Gopkg.lock ./
# RUN dep ensure --vendor-only
COPY . ./

ARG FLAGS
RUN CGO_ENABLED=0 GOOS=linux go build $FLAGS -a -o /app .

FROM debian:buster
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./

EXPOSE 8080
#ENTRYPOINT ["./app"]
CMD ["./app"]
