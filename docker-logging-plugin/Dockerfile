FROM  golang:1.20 as build
WORKDIR /app
COPY . .
RUN go build --ldflags '-extldflags "-static"' -o /usr/bin/docker-log-plugin

FROM alpine:3.18
COPY --from=build /usr/bin/docker-log-plugin /usr/bin/
WORKDIR /usr/bin
ENTRYPOINT [ "/usr/bin/docker-log-plugin" ]

#FROM cgr.dev/chainguard/static:latest
#COPY --from=build /usr/bin/docker-log-plugin /usr/bin/
#CMD ["/usr/bin/docker-log-plugin"]
