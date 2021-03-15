FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app ./cmd

FROM scratch
COPY --from=builder /app/app /usr/bin/app
ENTRYPOINT ["/usr/bin/app"]