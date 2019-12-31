FROM golang:1.13.5-alpine AS goBuilder
WORKDIR /builder
ADD backend/ .
#RUN apk update && apk add --no-cache git
#RUN go version
RUN go build -o lgn

FROM alpine
WORKDIR /app
COPY --from=goBuilder /builder/lgn /app/
COPY backend/migrations/ migrations/
ENTRYPOINT ["./lgn"]
