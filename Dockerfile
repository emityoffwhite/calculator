FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/calc ./cmd/calc

FROM alpine:3.20

RUN adduser -D -g '' appuser

COPY --from=builder /bin/calc /bin/calc

USER appuser

ENTRYPOINT ["/bin/calc"]
