FROM golang:1.21.13 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /kubeeye ./cmd/kubeeye

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /kubeeye /kubeeye
ENTRYPOINT ["/kubeeye"]
