FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/auction cmd/auction/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auction .

EXPOSE 8080

ENTRYPOINT ["./auction"]
