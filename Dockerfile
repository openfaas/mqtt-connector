FROM golang:1.12 as build

ENV GO111MODULE=off
ENV CGO_ENABLED=0

RUN mkdir -p /go/src/github.com/openfaas-incubator/mqtt-connector
WORKDIR /go/src/github.com/openfaas-incubator/mqtt-connector

COPY vendor     vendor
COPY main.go    .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"

RUN go test -v ./...

# Stripping via -ldflags "-s -w" 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /usr/bin/connector

FROM alpine:3.10 as ship

RUN apk add --no-cache ca-certificates

COPY --from=build /usr/bin/connector /usr/bin/connector
WORKDIR /root/

CMD ["/usr/bin/connector"]
