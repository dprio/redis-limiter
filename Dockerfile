FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o redis-limiter ./cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/redis-limiter .
EXPOSE 8080
EXPOSE 8081
EXPOSE 50051
CMD ["./redis-limiter"]