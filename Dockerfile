FROM golang:latest AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o lb .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/lb .
ENTRYPOINT [ "/root/lb" ]
