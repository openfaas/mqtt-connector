FROM ghcr.io/openfaas/license-check:0.4.0 as license-check

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18 as build

ARG GIT_COMMIT
ARG VERSION

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor


COPY --from=license-check /license-check /usr/bin/

WORKDIR /go/src/github.com/openfaas/mqtt-connector
COPY . .

RUN license-check -path /go/src/github.com/openfaas/mqtt-connector/ --verbose=false "Alex Ellis" "OpenFaaS Author(s)"
RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*")
RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go test -v ./...

RUN echo ${GIT_COMMIT} ${VERSION}
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=${CGO_ENABLED} go build \
        -mod=vendor \
        --ldflags "-s -w -X 'github.com/openfaas/mqtt-connector/version.GitCommit=${GIT_COMMIT}' -X 'github.com/openfaas/mqtt-connector/version.Version=${VERSION}'" \
        -a -installsuffix cgo -o mqtt-connector

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.14 as ship
LABEL org.label-schema.license="MIT" \
      org.label-schema.vcs-url="https://github.com/openfaas/mqtt-connector" \
      org.label-schema.vcs-type="Git" \
      org.label-schema.name="openfaas/mqtt-connector-pro" \
      org.label-schema.vendor="openfaas" \
      org.label-schema.docker.schema-version="1.0"

RUN apk --no-cache add \
    ca-certificates

RUN addgroup -S app \
    && adduser -S -g app app

WORKDIR /home/app

ENV http_proxy      ""
ENV https_proxy     ""

COPY --from=build /go/src/github.com/openfaas/mqtt-connector/mqtt-connector    /usr/bin/
RUN chown -R app:app ./

USER app

CMD ["/usr/bin/mqtt-connector"]
