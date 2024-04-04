FROM golang:1.21 as builder

WORKDIR /src
COPY . .
RUN go mod download
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o ./server/app -a ./cmd/main.go

FROM ubuntu:22.04
COPY --from=builder /src/server /server
COPY --from=builder /src/cmd/migrations /src/migrations
RUN apt-get update && apt-get install -y build-essential
RUN apt-get install -y ca-certificates
RUN ["chmod", "+x", "/server"]
EXPOSE 8000
CMD ["/server/app"]
