FROM alpine:latest
EXPOSE 8080
COPY main.go .
CMD go run main.go