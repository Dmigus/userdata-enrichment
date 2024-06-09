FROM golang:1.22-alpine AS build-stage


WORKDIR /app
COPY enrichstorage/go.mod go.mod
COPY enrichstorage/go.sum go.sum

RUN go mod download

COPY enrichstorage/ .

RUN go build -o bin/app ./cmd

FROM alpine AS release-stage

WORKDIR /app

COPY --from=build-stage /app/bin/app /app/bin/app

EXPOSE 50051
EXPOSE 8081

CMD ["./bin/app"]