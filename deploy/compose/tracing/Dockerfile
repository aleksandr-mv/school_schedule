FROM otel/opentelemetry-collector-contrib:0.133.0 AS base

FROM alpine:3.18
COPY --from=base /otelcol-contrib /otelcol-contrib

USER root
COPY otel-config.yaml /etc/otel-collector-config.yaml
ENTRYPOINT ["/otelcol-contrib", "--config", "/etc/otel-collector-config.yaml"]

EXPOSE 13133
