FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY main.go go.mod ./
RUN CGO_ENABLED=0 GOOS=linux go build -o balancer .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/balancer .
ENTRYPOINT [ "/root/balancer" ]