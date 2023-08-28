FROM golang:1.21 as build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o admission-webhook-demo ./

FROM cgr.dev/chainguard/static:latest
LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
COPY --from=build /app/admission-webhook-demo /admission-webhook-demo
ENTRYPOINT ["/admission-webhook-demo"]
