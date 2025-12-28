FROM golang:1.25 AS build
USER root
WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 go build -o delay-proxy .

FROM cgr.dev/chainguard/static:latest
LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
COPY --from=build /work/delay-proxy /delay-proxy
ENTRYPOINT ["/delay-proxy"]
