## 1. 构建应用 ##
FROM golang:1.15-buster as builder

ARG APP_MODULE="github.com/Gopusher/gateway"
ARG APP_NAME="gateway"

RUN export GO111MODULE=on \
    && export GOPROXY=https://goproxy.io \
    && mkdir -p /go/src/${APP_MODULE}

COPY . /go/src/${APP_MODULE}

# RUN echo /go/src/${APP_MODULE}
# RUN echo /go/bin/${APP_NAME} ${APP_NAME}.go

RUN cd /go/src/${APP_MODULE} \
    && CGO_ENABLED=0 go build -ldflags '-s -w' -o /go/bin/${APP_NAME} ./app/${APP_NAME}/app/cmd/main.go

## 2. 应用 ##
FROM debian:buster

ARG OLD_APP_NAME="gateway"
ARG APP_NAME="gopusher"


# 增加根证书
RUN apt-get update \
    && apt-get install ca-certificates -y

COPY --from=builder /go/bin/${OLD_APP_NAME} /usr/local/bin/${APP_NAME}

# .env => /app/.env
WORKDIR /app
VOLUME /app

EXPOSE 8900 8901

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["${APP_NAME}", "start", "-c", "config.yaml"]
