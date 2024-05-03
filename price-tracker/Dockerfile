FROM golang:1.22 as build
USER root
WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 go build -o price-tracker .

FROM cgr.dev/chainguard/static:latest
LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
COPY --from=build /work/price-tracker /price-tracker
ENTRYPOINT ["/price-tracker"]
