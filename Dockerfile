FROM gcr.io/distroless/base

COPY bin/straw-pod /straw-pod

ENTRYPOINT ["/straw-pod"]