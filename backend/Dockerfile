FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=build /api /app/api
COPY migrations ./migrations
COPY seeds ./seeds

EXPOSE 8080

ENTRYPOINT ["/app/api"]