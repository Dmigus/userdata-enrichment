FROM golang:1.22-alpine AS build-stage

WORKDIR /app/

COPY go.work go.work
COPY go.work.sum go.work.sum

COPY enrichstorage/pkg enrichstorage/pkg
COPY enrichstorage/go.mod enrichstorage/go.mod
COPY enrichstorage/go.sum enrichstorage/go.sum

WORKDIR /app/enricher
COPY enricher/go.mod go.mod
COPY enricher/go.sum go.sum

RUN go mod download

COPY ./enricher .

RUN go build -o bin/app ./cmd

FROM alpine AS release-stage

WORKDIR /app

COPY --from=build-stage /app/enricher/bin/app /app/bin/app

EXPOSE 8080

CMD ["./bin/app"]