FROM golang:1.25 AS build
USER root
WORKDIR /work
COPY . /work
RUN CGO_ENABLED=0 go build -o proxy .

FROM gcr.io/distroless/static-debian13:debug
LABEL org.opencontainers.image.source="https://github.com/morawskim/go-projects"
COPY --from=build /work/proxy /proxy
USER nonroot
EXPOSE 8081
WORKDIR /home/nonroot
ENTRYPOINT ["/proxy"]
