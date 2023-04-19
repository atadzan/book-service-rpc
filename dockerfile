FROM golang:1.19.5-alpine3.17

RUN echo 'https://dl-3.alpinelinux.org/alpine/v3.17/main' > /etc/apk/repositories  && \
    echo '@testing https://dl-3.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    echo '@community https://dl-3.alpinelinux.org/alpine/v3.17/community'

RUN apk update && apk add --no-cache git
RUN apk add postgresql-client

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o book-service-rpc ./cmd/server/main.go

CMD ["./book-service-rpc"]