# This is the base specification for building master agent as a docker image.
# Use [tools/build-docker] to easily build the image.
# Use [tools/run-docker] to easily run the container.

FROM golang:1.23.1-alpine AS builder
WORKDIR /zero-monitor

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY public/ ./public/
COPY .version .version

RUN version="$(cat .version)"
RUN build_pkg="github.com/guackamolly/zero-monitor/internal/build"
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-X $build_pkg.version=$version -X $build_pkg.release=true -extldflags=-static -s -w" -o master ./cmd/master/master.go

FROM scratch
WORKDIR /zero-monitor

COPY --from=builder /zero-monitor /zero-monitor
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "/zero-monitor/master" ]
