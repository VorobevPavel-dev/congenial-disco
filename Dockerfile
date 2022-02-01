FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/disco

FROM scratch
WORKDIR /app
EXPOSE 8090
COPY --from=builder /app/disco /usr/bin/
ENTRYPOINT ["disco"]