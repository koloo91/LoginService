FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /src
ADD . .
RUN go build -o lgn

FROM alpine
WORKDIR /app/
COPY --from=builder /src/ .
ENTRYPOINT ["/app/lgn"]
