FROM golang:alpine AS builder
WORKDIR /src
COPY . /src
RUN apk add --no-cache git
RUN go test -c -o /tmp/e2e.test ./e2e/

FROM alpine:latest
COPY --from=builder /tmp/e2e.test .
RUN apk add --no-cache ca-certificates
ENTRYPOINT ["./e2e.test", "-test.v", "-test.timeout=10m"]
