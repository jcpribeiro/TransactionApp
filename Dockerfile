ARG GOLANG_VERSION


FROM golang:${GOLANG_VERSION}-alpine as builder
WORKDIR /build
RUN apk update --no-cache \
    && apk add --no-cache build-base git openssh \
    && rm -rf /var/cache/apk/*

ADD go.mod .
ADD go.sum .

RUN go mod download
COPY . .
RUN env GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o /bin/app main.go



FROM alpine:3.16 as application
RUN apk update --no-cache \
    && apk add -U --no-cache ca-certificates tzdata \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /bin/app .
COPY --from=builder /build/config.json .
COPY --from=builder /build/config_prod.json .

CMD ["/app/app"]