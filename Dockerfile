FROM golang:alpine AS build_base

RUN apk add --no-cache git && \
    apk add openssl

WORKDIR /tmp/chat

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/app .

FROM alpine:latest

ENV GIN_MODE=release

COPY --from=build_base /tmp/chat/out/app /app/chat

COPY public /app/public
COPY config.yml /app

WORKDIR /app

ENV DB_PASSWORD postgres
ENV SECRET_KEY 4r7hud9872jfulpqcb40HycaGf63nZX950p1kdmqnaThqlKVysneRpl91Qejkie

EXPOSE 8080

ENTRYPOINT ["/app/chat"]