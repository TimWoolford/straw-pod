FROM gcr.io/distroless/base

COPY web/html /html

COPY bin/straw-pod /straw-pod

ENTRYPOINT ["/straw-pod"]