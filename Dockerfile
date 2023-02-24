FROM golang:alpine AS builder
WORKDIR /app
COPY go.* ./
COPY . ./
RUN go mod download && go build -tags=jsoniter -o emenu-api cmd/main.go

FROM alpine:3.17.0
WORKDIR /app
RUN addgroup -S user && adduser -S user -G user

COPY --from=builder /app/emenu-api /app/emenu-api
COPY .env  /app/

RUN chown -R user:user /app
USER user

CMD ["/app/emenu-api"]