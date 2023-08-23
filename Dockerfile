FROM golang:1.19.12-alpine3.18 AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY exporter.go ./
RUN go mod download
RUN go build -o app

FROM --platform=linux/amd64 alpine:latest
WORKDIR /
COPY --from=build /app/app .
#RUN chmod +x /app/app
EXPOSE 9100
ENTRYPOINT ["/app"]