FROM golang:1.20

ADD ../.. /app

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache go build -o /api cmd/url-shortener/main.go

EXPOSE 8080

CMD /api -in-memory

ENV CONFIG_PATH=./config/local.yaml