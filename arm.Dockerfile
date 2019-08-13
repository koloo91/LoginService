FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /src
ADD . .
RUN env GOOS=linux GOARCH=arm go build -o lgn

FROM armhf/alpine:latest
WORKDIR /app/
COPY --from=builder /src/ .
ENTRYPOINT ["/app/lgn"]
