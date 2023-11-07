FROM golang:1.21 as build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o tapo-metrics ./

FROM cgr.dev/chainguard/static:latest
LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
COPY --from=build /app/tapo-metrics /tapo-metrics
ENTRYPOINT ["/tapo-metrics"]
