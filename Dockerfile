FROM golang:1.15.3-alpine as builder
LABEL maintainer="Brian Marin"

RUN apk update && apk add --virtual build-dependencies build-base ca-certificates

WORKDIR /build
COPY . .

RUN go mod download

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags="-w -s" -o something_server cmd/something/backend/main.go

FROM scratch

WORKDIR /opt/something

COPY --from=builder /build/something_server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT [ "./something_server" ]