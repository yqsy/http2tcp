FROM golang:1.12.5-alpine3.9

RUN set -ex; \
    echo "https://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories; \
    apk update; \
    apk add bash; \
    mkdir -p /tmp/http2tcp;


COPY ./latest.tar.gz /tmp/http2tcp

RUN set -ex; \
    cd /tmp/http2tcp; \
    tar -xvzf latest.tar.gz; \
    go build -o http2tcp main.go; \
    mv http2tcp /usr/local/bin;

# 50003 http 接口, 转50001的tcp line接口
# EXPOSE 50003