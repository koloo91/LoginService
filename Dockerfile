FROM golang:1.13.4-alpine AS goBuilder
WORKDIR /builder
ADD backend .
RUN go build -o lgn

FROM alpine
WORKDIR /app
COPY --from=goBuilder /builder/lgn /app/
COPY backend/migrations/ migrations/
ENTRYPOINT ["./lgn"]
