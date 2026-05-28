FROM golang:1.21.4 AS builder
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY pkg ./pkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-binary ./cmd/kubeeye

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app-binary /app-binary
EXPOSE 8888
CMD ["/app-binary"]
