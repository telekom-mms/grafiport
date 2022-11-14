# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS build-env
RUN mkdir -p /go/src/grafana-exporter

# Copy the module files first and then download the dependencies. If this
# doesn't change, we won't need to do this again in future builds.
COPY go.* /go/src/grafana-exporter/
WORKDIR /go/src/grafana-exporter
RUN go mod download

WORKDIR /go/src/grafana-exporter
ADD cmd cmd
COPY main.go /go/src/grafana-exporter/
WORKDIR /go/src/grafana-exporter
RUN go build -o grafana-exporter

# final stage
FROM alpine:latest
COPY --from=build-env /go/src/grafana-exporter/grafana-exporter /app/
RUN mkdir -p /opt/data/grafana-exporter
ENV DIRECTORY /opt/data/grafana-exporter
ENTRYPOINT /app/grafana-exporter