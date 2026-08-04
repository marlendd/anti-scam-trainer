FROM golang:1.25 AS build

COPY go.mod go.sum /src/
COPY cmd/api /src/api

RUN cd /src && go build -o /api api/main.go

FROM alpine:3.20

COPY --from=build /api /api

ENTRYPOINT [ "/api" ]