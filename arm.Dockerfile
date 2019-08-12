FROM arm32v7/golang:1.12.7-alpine3.10 AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /src
ADD . .
RUN go build -o lgn

FROM arm32v6/alpine
WORKDIR /app/
COPY --from=builder /src/ .
ENTRYPOINT ["/app/lgn"]
