####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/grafana-log-sink-${TARGETARCH} /bin/grafana-log-sink
RUN chmod +x /bin/grafana-log-sink

####################################################################################################
# grafana-log-sink
####################################################################################################
FROM scratch AS grafana-log-sink
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/grafana-log-sink /bin/grafana-log-sink
ENTRYPOINT [ "/bin/grafana-log-sink" ]
