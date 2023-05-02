# final stage
FROM alpine:latest
COPY grafana-exporter /usr/local/bin/grafana-exporter
RUN mkdir -p /output
ENV DIRECTORY /output
ENTRYPOINT ["grafana-exporter"]
