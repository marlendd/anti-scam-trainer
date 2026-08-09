FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /api ./cmd/api

FROM alpine:3.20

COPY --from=build /api /api

ENTRYPOINT [ "/api" ]