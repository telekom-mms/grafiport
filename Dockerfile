# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build-env
RUN mkdir -p /go/src/grafiport

# Copy the module files first and then download the dependencies. If this
# doesn't change, we won't need to do this again in future builds.
COPY go.* /go/src/grafiport/
WORKDIR /go/src/grafiport
RUN go mod download

WORKDIR /go/src/grafiport
ADD export export
ADD restore restore
ADD common common
COPY *.go /go/src/grafiport/
WORKDIR /go/src/grafiport
RUN go build -o grafiport

# final stage
FROM alpine:latest
COPY --from=build-env /go/src/grafiport/grafiport  /usr/local/bin/grafiport
RUN mkdir -p /output
ENV DIRECTORY /output
ENTRYPOINT ["grafiport"]