####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/grafana-sink-${TARGETARCH} /bin/grafana-sink
RUN chmod +x /bin/grafana-sink

####################################################################################################
# grafana-sink
####################################################################################################
FROM scratch AS grafana-sink
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/grafana-sink /bin/grafana-sink
ENTRYPOINT [ "/bin/grafana-sink" ]
