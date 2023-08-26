FROM gcr.io/distroless/static-debian11
ARG TARGETARCH

LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
ENTRYPOINT ["/trigger-lambda"]
COPY trigger-lambda-$TARGETARCH /trigger-lambda
