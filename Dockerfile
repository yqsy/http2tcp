FROM golang:1.12.5-alpine3.9

RUN set -ex; \
    echo "https://mirrors.aliyun.com/alpine/v3.9/main/" > /etc/apk/repositories; \
    apk update; \
    apk add --no-cache bash apache2-utils;

COPY . /env/http2tcp

RUN set -ex; \
    cd /env/http2tcp; \
    go build -o http2tcp main.go; \
    mv http2tcp /usr/local/bin; \
    cd /env; \
    rm -rf /env/http2tcp

# 50003 http 接口, 转50001的tcp line接口
EXPOSE 50003