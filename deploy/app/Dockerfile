FROM golang:1.20

ADD ../.. /app

ENV CONFIG_PATH=./config/local.yaml
ARG PROTOCOL=cmd/url-shortener/main.go

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache  \
    --mount=type=cache,target=/go-cache \
    go build -o /api ${PROTOCOL}

EXPOSE 8080

CMD /api -in-memory

