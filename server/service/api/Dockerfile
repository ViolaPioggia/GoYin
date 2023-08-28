FROM alpine:latest

WORKDIR /

COPY ./api .
COPY ./config.yaml ./server/service/api/config.yaml
COPY ./IP2LOCATION-LITE-DB11.BIN ./IP2LOCATION-LITE-DB11.BIN

# 暴露端口
EXPOSE 8080

CMD ["./api"]
