FROM golang:1.24 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cloudrun

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]