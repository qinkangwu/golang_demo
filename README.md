## dockerfile build 命令
```shell
DOMAIN=$1
DOCKERIMAGENAME=$2
cd ../server2
docker build -t $DOCKERIMAGENAME -f ../dev/$DOMAIN/Dockerfile .
```

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

```shell
docker run -d --restart always --name prometheus -p 9190:9090 -v /Users/qinkangwu/goStudy/prometheus.yml:/etc/prometheus/prometheus.yml -v /Users/qinkangwu/goStudy/rules/target.yml:/etc/prometheus/rules/target.yml  prom/prometheus


docker run --name alertmanager -d -p 0.0.0.0:9093:9093 quay.io/prometheus/alertmanager 

docker run -d -p 3000:3000 --name=grafana -v /Users/qinkangwu/prometheus-data/grafana-storage/:/var/lib/grafana grafana/grafana  

docker run -d -p 9100:9100 -v "/proc:/host/proc:ro" -v "/sys:/host/sys:ro" -v "/:/rootfs:ro" prom/node-exporter  
```