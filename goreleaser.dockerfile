# final stage
FROM alpine:latest
COPY grafiport /usr/local/bin/grafiport
RUN mkdir -p /output
ENV DIRECTORY /output
ENTRYPOINT ["grafiport"]
