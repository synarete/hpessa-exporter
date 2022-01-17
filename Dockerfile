FROM golang:1.16 AS builder

# args
ARG ARCH
ARG VERSION
ARG GITREV
ARG GOVARS
RUN echo $ARCH
RUN echo $VERSION
RUN echo $GITREV
RUN echo $GOVARS

# prepare
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/ cmd/
COPY internal/ internal/
COPY vendor/ vendor/

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH \
  go build -a -o /workspace/hpessa-exporter \
  -ldflags="-s -w $GOVARS" cmd/main.go

# create distroless image
FROM gcr.io/distroless/static:nonroot
ARG ARCH
WORKDIR /
COPY --from=builder /workspace/hpessa-exporter .
USER 65532:65532

ENTRYPOINT ["/hpessa-exporter"]
EXPOSE 8080
