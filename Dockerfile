FROM golang:alpine AS builder
WORKDIR /app
COPY go.* ./
COPY . ./
RUN go mod download
RUN go build -o emenu-api cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/emenu-api /app/emenu-api
COPY .env  /app/
CMD ["/app/emenu-api"]