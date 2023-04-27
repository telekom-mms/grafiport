# final stage
FROM scratch
COPY grafana-exporter /
ENTRYPOINT ["/grafana-exporter"]
