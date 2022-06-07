## dockerfile build 命令
``docker build -t qkw -f ../dev/gateway/Dockerfile .``

## dockerfile 内容
```dockerfile
FROM golang:1.18.2-alpine AS builder

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . /go/src/Demo/server2

WORKDIR /go/src/Demo/server2
RUN CGO_ENABLED=0 GOOS=linux go install ./rental/...

FROM alpine:3.15

COPY --from=builder /go/bin/rental /bin/rental
ENV ADDR=:8081
EXPOSE 8081

ENTRYPOINT ["/bin/rental"]
```