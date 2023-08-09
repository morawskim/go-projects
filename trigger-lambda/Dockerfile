FROM cgr.dev/chainguard/static:latest
ARG TARGETARCH

LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
ENTRYPOINT ["/trigger-lambda"]
COPY trigger-lambda-$TARGETARCH /trigger-lambda
